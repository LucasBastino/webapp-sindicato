package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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
	query := "SELECT * FROM payments WHERE year = ? AND id_company = ? ORDER BY month ASC";
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

type OverdueSummary struct {
	Count  int     `db:"count"`
	Amount float32 `db:"amount"`
}

func (r *PaymentRepository) CountOverdueSummary(ctx context.Context) (OverdueSummary, error) {
	query := `
		SELECT COUNT(*) AS count, COALESCE(SUM(amount), 0) AS amount
		FROM payments
		WHERE paid_at IS NULL
		  AND is_in_payment_plan = false
		  AND due_date < CURDATE()
	`
	var summary OverdueSummary
	err := r.db.GetContext(ctx, &summary, query)
	if err != nil {
		return OverdueSummary{}, fmt.Errorf("failed to count overdue payments: %w", err)
	}
	return summary, nil
}

type overdueRow struct {
	ID          int       `db:"id_payment"`
	CompanyID   int       `db:"id_company"`
	CompanyName string    `db:"company_name"`
	Month       int       `db:"month"`
	Year        int       `db:"year"`
	DueDate     time.Time `db:"due_date"`
	Amount      *float32  `db:"amount"`
}

func (r *PaymentRepository) FindAllOverdue(ctx context.Context) ([]overdueRow, error) {
	query := `
		SELECT
			P.id_payment,
			P.id_company,
			C.name AS company_name,
			P.month,
			P.year,
			P.due_date,
			P.amount
		FROM payments P
		INNER JOIN companies C ON P.id_company = C.id_company
		WHERE P.paid_at IS NULL
		  AND P.is_in_payment_plan = false
		  AND P.due_date < CURDATE()
		  AND C.deleted_at IS NULL
		ORDER BY C.name ASC, P.due_date ASC
	`
	var rows []overdueRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch overdue payments: %w", err)
	}
	return rows, nil
}

func (r *PaymentRepository) FindOverdueByCompany(ctx context.Context, companyID int) ([]Payment, error) {
	query := `
		SELECT *
		FROM payments
		WHERE id_company = ?
		  AND paid_at IS NULL
		  AND is_in_payment_plan = false
		  AND due_date < CURDATE()
		ORDER BY due_date ASC, month ASC
	`
	var payments []Payment
	err := r.db.SelectContext(ctx, &payments, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch overdue payments by company: %w", err)
	}
	return payments, nil
}

func (r *PaymentRepository) FindByIDs(ctx context.Context, ids []int) ([]Payment, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(`SELECT * FROM payments WHERE id_payment IN (?)`, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to build find-by-ids query: %w", err)
	}
	query = r.db.Rebind(query)
	var payments []Payment
	err = r.db.SelectContext(ctx, &payments, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments by ids: %w", err)
	}
	return payments, nil
}

func (r *PaymentRepository) BulkUpdateIsInPaymentPlan(ctx context.Context, tx *sqlx.Tx, ids []int, inPlan bool) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	query, args, err := sqlx.In(`UPDATE payments SET is_in_payment_plan = ? WHERE id_payment IN (?)`, inPlan, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to build bulk is_in_payment_plan query: %w", err)
	}
	query = r.db.Rebind(query)
	var res sql.Result
	if tx != nil {
		res, err = tx.ExecContext(ctx, query, args...)
	} else {
		res, err = r.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to bulk update is_in_payment_plan: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected while bulk updating is_in_payment_plan: %w", err)
	}
	return int(rows), nil
}

func (r *PaymentRepository) BulkMarkCompleted(ctx context.Context, tx *sqlx.Tx, ids []int, paidAt time.Time) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	query, args, err := sqlx.In(`
		UPDATE payments
		SET paid_at = ?, is_in_payment_plan = false
		WHERE id_payment IN (?)
	`, paidAt, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to build bulk mark completed query: %w", err)
	}
	query = r.db.Rebind(query)
	var res sql.Result
	if tx != nil {
		res, err = tx.ExecContext(ctx, query, args...)
	} else {
		res, err = r.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to bulk mark payments completed: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected while bulk marking payments completed: %w", err)
	}
	return int(rows), nil
}


// Use INSERT IGNORE to skip duplicates payments
// Also skip softdeleted companies with condition WHERE 
func (r *PaymentRepository) BulkInsert(ctx context.Context, tx *sqlx.Tx, payments []Payment) error {
	query := "INSERT IGNORE INTO payments (id_company, month, year, due_date, observations) VALUES"
	placeholders := make([]string, 0, len(payments))
	args := make([]any, 0, len(payments)*5)
	for _, p := range payments{
		placeholders = append(placeholders, " (?, ?, ?, ?, ?)")
		args = append(args, p.CompanyID, p.Month, p.Year, p.DueDate, "")
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
	return nil
}


