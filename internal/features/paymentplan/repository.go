package paymentplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PaymentPlanRepository struct {
	db *sqlx.DB
}

func NewPaymentPlanRepository(db *sqlx.DB) *PaymentPlanRepository {
	return &PaymentPlanRepository{
		db: db,
	}
}

func (r *PaymentPlanRepository) FindByID(ctx context.Context, id int) (*PaymentPlan, error) {
	query := "SELECT * FROM payment_plans WHERE id_payment = ?"
	var paymentPlan PaymentPlan
	err := r.db.GetContext(ctx, &paymentPlan, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch payment plan: %w", err)
	}
	return &paymentPlan, nil
}

func (r *PaymentPlanRepository) FindAll(ctx context.Context, companyID int) ([]PaymentPlan, error) {
	query := `
	SELECT * FROM payment_plans P
	INNER JOIN companies C
		ON P.id_company = C.id_company
	WHERE
		P.id_company = ?
	`
	var paymentPlans []PaymentPlan
	err := r.db.GetContext(ctx, paymentPlans, query, companyID)
	if err!=nil{
		return nil, fmt.Errorf("failed to fetch payment plans: %w", err)
	}
	return paymentPlans, nil	
}

func (r *PaymentPlanRepository) Insert(ctx context.Context, paymentPlan PaymentPlan) (int, error) {
	query := 
		`INSERT INTO payment_plans (
			id_company,
			amount,
			number_of_installments,
			first_due_date,
			observations
		) VALUES (
			:id_company,
			:amount,
			:number_of_installments,
			:observations
		)`


	res, err := r.db.NamedExecContext(ctx, query, paymentPlan)
	if err!=nil{
		return 0, fmt.Errorf("failed to insert payment plan: %w", err)
	}

	id, err := res.LastInsertId()
	if err!=nil{
		return 0, fmt.Errorf("failed to get last insert id while inserting payment plan: %w", err)
	}
	return int(id), nil
}

func (r *PaymentPlanRepository) Update(ctx context.Context, id int, paymentPlan PaymentPlan) error {
	query := "UPDATE payment_plans SET observations = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, paymentPlan.Observations, id)
	if err!=nil{
		return fmt.Errorf("failed to restore payment plan: %w", err)
	}
	return nil
}

func (r *PaymentPlanRepository) Cancel(ctx context.Context, tx *sqlx.Tx, id int) (int, error) {
	query := "UPDATE payment_plans SET status = 'cancelled' WHERE id = ?"
	res, err := tx.ExecContext(ctx, query, id)
	if err!=nil{
		return 0, fmt.Errorf("failed to cancel payment plan: %w", err)
	}

	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while cancelling payment plan: %w", err)
	}
	return int(rows), nil
}

func (r *PaymentPlanRepository) Restore(ctx context.Context, tx *sqlx.Tx, id int) (int, error) {
	query := "UPDATE payment_plans SET status = 'pending' WHERE id = ?"
	res, err := tx.ExecContext(ctx, query, id)
	if err!=nil{
		return 0, fmt.Errorf("failed to restore payment plan: %w", err)
	}

	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while restoring payment plan: %w", err)
	}
	return int(rows), nil
}

func (r *PaymentPlanRepository) UpdateStatus(ctx context.Context, tx *sqlx.Tx, id int, status string) error {
	query := "UPDATE payment_plans SET status = ? WHERE id = ?"
	_, err := tx.ExecContext(ctx, query, status, id)
	if err!=nil{
		return fmt.Errorf("failed to update payment plan status: %w", err)
	}
	return nil
}

func (r *PaymentPlanRepository) HardDelete(ctx context.Context, id int) error {
	query := "DELETE FROM payment_plans WHERE id_payment_plan = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete payment plan: %w", err) 
	}
	return nil
}

func (r *PaymentPlanRepository) Count(ctx context.Context, companyID int) (int, error) {
	query := "SELECT COUNT(*) FROM payment_plans WHERE id_company = ?"
	var totalRows int
	err := r.db.GetContext(ctx, &totalRows, query, companyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total rows while fetching payment plans: %w", err) 
	}
	return totalRows, nil
}

func (r *PaymentPlanRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error){
	return r.db.BeginTxx(ctx, nil)
}