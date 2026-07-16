package installment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type InstallmentRepository struct {
	db *sqlx.DB
}

func NewInstallmentRepository(db *sqlx.DB) *InstallmentRepository{
	return &InstallmentRepository{
		db: db,
	}
}

func (r *InstallmentRepository) FindByID(ctx context.Context, id int) (*Installment, error) {
	query := "SELECT * FROM installments WHERE id_installment = ?"
	var installment Installment
	err := r.db.GetContext(ctx, &installment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch installment: %w", err)
	}
	return &installment, nil
}

func (r *InstallmentRepository) FindAll(ctx context.Context, tx *sqlx.Tx, paymentPlanID int) ([]Installment, error) {
	query := `
	SELECT * FROM installments I
	INNER JOIN payment_plans P
		ON I.id_payment_plan = P.id_payment_plan
	WHERE
		I.id_payment_plan = ?
	ORDER BY number_of_installment ASC
	`
	var installments []Installment
	var err error
	if tx != nil{
		err = tx.GetContext(ctx, installments, query, paymentPlanID)	
	} else{
		err = r.db.GetContext(ctx, installments, query, paymentPlanID)
	}
	if err!=nil{
		return nil, fmt.Errorf("failed to fetch installments: %w", err)
	}
	return installments, nil
}

func (r *InstallmentRepository) BulkInsert(ctx context.Context, installments []Installment) (int, error) {
	query := `
	INSERT INTO installments (
		id_payment_plan,
		installment_number,
		amount,
		due_date
	) VALUES
	`
	// preallocate placeholders and args
	placeholders := make([]string, 0, len(installments))
	args := make([]any, 0, len(installments))

	for _, i := range installments{
		placeholders = append(placeholders, " (?, ?, ?, ?)")
		args = append(args, i.PaymentPlanID, i.InstallmentNumber, i.Amount, i.DueDate)
	}

	query += strings.Join(placeholders, ",")

	res, err := r.db.ExecContext(ctx, query, args...)
	if err!=nil{
		return 0, fmt.Errorf("failed to bulk insert installments: %w", err)
	}
	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while bulk inserting installments: %w", err)
	}
    return int(rows), nil
}

func (r *InstallmentRepository) Update(ctx context.Context, id int, installment Installment) error{
	installment.ID = id

	query := `
	UPDATE installments
	SET
		status = :status,
		observations = :observations
	WHERE 
		id_installment = :id_installment
	`
	
	_, err := r.db.NamedExecContext(ctx, query, installment)
	if err != nil {
		return fmt.Errorf("failed to update installment: %w", err) 
	}
	return nil
}

func (r *InstallmentRepository) CheckStatusDaily(ctx context.Context) error {
	query := `
	UPDATE installments
	SET
		status = 'overdue'
	WHERE
		status = 'pending'
	AND
		due_date < CURDATE()
	`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute daily installments update: %w", err) 
	}
	return nil
}




