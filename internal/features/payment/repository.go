package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository struct{
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository{
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) FindByID(ctx context.Context, id int) (*Payment, error) {
	query := "SELECT * FROM payments WHERE id_payment = ?";
	var payment Payment
	err := r.db.GetContext(ctx, &payment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch payment: %w", err)
	}
	return &payment, nil
}

func (r *PaymentRepository) FindAll(ctx context.Context, companyID int, year int) ([]Payment, error) {
	query := "SELECT * FROM payments WHERE year = ? AND id_company = ? ORDER BY month DESC";
	var payments []Payment
	err := r.db.SelectContext(ctx, &payments, query, year, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}
	return payments, nil
}

func (r *PaymentRepository) GetPaymentYears(ctx context.Context, companyID int) ([]int, error){
	query := "SELECT year FROM payments WHERE id_company = ? GROUP BY year ORDER BY year DESC"
	var years []int
	err := r.db.SelectContext(ctx, &years, query, companyID)
	if err != nil {
		return nil,  fmt.Errorf("failed to fetch payment years: %w", err)
	}
	return years, nil
}

func (r *PaymentRepository) Count(ctx context.Context, companyID int) (int, error){
	query := "SELECT COUNT(*) FROM payments WHERE id_company = ?"
	var totalRows int
	err := r.db.GetContext(ctx, &totalRows, query, companyID)
	if err != nil {
		return 0,  fmt.Errorf("failed to count payments: %w", err)
	}
	return totalRows, nil
}


// Use INSERT IGNORE to skip duplicates payments
// Also skip softdeleted companies with condition WHERE 
func (r *PaymentRepository) BulkInsert(ctx context.Context, tx *sqlx.Tx, payments []Payment) error {
	query := "INSERT IGNORE INTO payments (id_company, month, year) VALUES"
	placeholders := make([]string, 0, len(payments))
	args := make([]any, 0, len(payments))
	for _, p := range payments{
		placeholders = append(placeholders, " (?, ?, ?)")
		args = append(args, p.CompanyID, p.Month, p.Year)
	}

	query += strings.Join(placeholders, ",")
	var err error
	if tx!=nil{
		_, err = tx.ExecContext(ctx, query, args...)
	} else{
		_, err = r.db.ExecContext(ctx, query, args...)
	}
    if err != nil {
        return fmt.Errorf("failed to bulk insert payments: %w", err)
    }
    return nil
}

// no se puede editar mes, año ni id_company
// no se puede editar si la empresa esta softdeleted
func (r *PaymentRepository) Update(ctx context.Context, id int, payment Payment) error {
	payment.ID = id
	query := `
		UPDATE payments P
		INNER JOIN companies C ON P.id_company = C.id_company
		SET 
			P.status = :status, 
			P.amount = :amount, 
			P.paid_at = :paid_at, 
			P.observations = :observations 
		WHERE 
			P.id_payment = :id_payment
			AND C.deleted_at IS NULL `;
	_, err := r.db.NamedExecContext(ctx, query, payment)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	return  nil
}

func (r *PaymentRepository) BulkUpdateStatus(ctx context.Context, tx *sqlx.Tx, ids string, status string) (int, error) {
	
	query := `
	UPDATE installments
	SET status = ?
	WHERE id IN 
	`
	// se que la cantidad de "?" va a ser exactamente la cantidad de ids, asi que ya lo voy creando
	placeholders := make([]string, len(ids))
	// lo mismo para los args, solo que se le suma el arg status, por eso +1
	args := make([]any, 0, len(ids)+1)

	args = append(args, status)

	for i, id := range ids {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query+= " (" + strings.Join(placeholders, ",") + ")"

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to bulk update installments: %w", err)
	}
	
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected while bulk updating installments: %w", err)
	}
	return int(rows), nil

}

func (r *PaymentRepository) CheckStatusMonthly(ctx context.Context) error {
	query := `
	UPDATE payments
	SET
		status = 'overdue'
	WHERE
		status = 'pending'
	AND
		DAY(CURDATE()) > 15
	`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute monthly payments update: %w", err)
	}
	return nil
}

// func (r *PaymentRepository) Insert(ctx context.Context, payment Payment) (int64, error) {
// 	query := `
// 	INSERT INTO payment 
// 	(month,
// 	year,
// 	status, 
// 	amount, 
// 	paid_at, 
// 	observations,
// 	id_company)
// 	VALUES (:month, :year, :status, :amount, :paid_at, :observations, :id_company)`;
// 	res, err := r.db.NamedExecContext(ctx, query, payment)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to insert payment: %w", err)
// 	}
// 	id, err := res.LastInsertId()
// 	if err!=nil{
// 		return 0, fmt.Errorf("failed to get last insert id of payment: %w", err)
// 	}
// 	return id, nil
// }



// func (r *PaymentRepository) HardDeleteAllPaymentsByCompany(ctx context.Context, companyID int) (int64, error) {
// 	query := "DELETE FROM payments WHERE id_company = ?"
// 	res, err := r.db.ExecContext(ctx, query, companyID)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to hard delete payments: %w", err)
// 	}
// 	rows, err := res.RowsAffected()
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get rows affected while hard deleting payments: %w", err)
// 	}
	
// 	return rows, nil
// }








// func (r *PaymentRepository) Exists(ctx context.Context, id int) (bool, error){
// 	query := "SELECT EXISTS (SELECT 1 FROM payments WHERE id_payment = ?)"
// 	var exists bool
// 	err := r.db.GetContext(ctx, &exists, query, id)
// 	if err!=nil{
// 		return false, fmt.Errorf("failed to check if member exists: %w", err)
// 	}
// 	return exists, nil
// }

// func (r *PaymentRepository) ValidateInDB(ctx context.Context, companyID int) map[string]string{
// 	errorMap := map[string]string{}
// 	if err := r.validateCompanyID(ctx, companyID); err!=""{
// 		errorMap["companyID"] = err
// 	}
// 	return errorMap
// }

// func (r *PaymentRepository) validateCompanyID(ctx context.Context, companyID int) string{

// 	var dummy int
// 	query := "SELECT 1 FROM companies WHERE id_company = ?";
// 	err := r.db.GetContext(ctx, &dummy, query, companyID)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return "Empresa no existente."
// 		}
// 		return "Ocurrió un error con el campo."

// 	}
// 	return ""
// }

// func (r *PaymentRepository) BeginTx() (*sqlx.Tx, error){
// 	return r.db.Beginx()
// }