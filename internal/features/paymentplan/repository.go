package paymentplan

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type PaymentPlanRepository struct {
	db *sqlx.DB
}

func NewPaymentPlanRepository(db *sqlx.DB) *PaymentPlanRepository {
	return &PaymentPlanRepository{db: db}
}

func (r *PaymentPlanRepository) FindByID(ctx context.Context, id int) (*PaymentPlan, error) {
	query := `
		SELECT
			PP.id_payment_plan,
			PP.id_company,
			C.name AS company_name,
			COALESCE(PP.payments_in_plan, '') AS payments_in_plan,
			PP.amount,
			PP.number_of_installments,
			COALESCE(PP.status, 'pending') AS status,
			PP.first_due_date,
			PP.last_due_date,
			COALESCE(PP.observations, '') AS observations,
			PP.created_at,
			PP.updated_at
		FROM payment_plans PP
		INNER JOIN companies C ON PP.id_company = C.id_company
		WHERE PP.id_payment_plan = ?
	`
	var paymentPlan PaymentPlan
	err := r.db.GetContext(ctx, &paymentPlan, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch payment plan: %w", err)
	}
	return &paymentPlan, nil
}

func (r *PaymentPlanRepository) FindAll(ctx context.Context, companyID int) ([]PaymentPlan, error) {
	query := `
		SELECT
			id_payment_plan,
			id_company,
			COALESCE(payments_in_plan, '') AS payments_in_plan,
			amount,
			number_of_installments,
			COALESCE(status, 'pending') AS status,
			first_due_date,
			last_due_date,
			COALESCE(observations, '') AS observations,
			created_at,
			updated_at
		FROM payment_plans
		WHERE id_company = ?
		ORDER BY updated_at DESC
	`
	var paymentPlans []PaymentPlan
	err := r.db.SelectContext(ctx, &paymentPlans, query, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment plans: %w", err)
	}
	return paymentPlans, nil
}

type planWithCompanyRow struct {
	ID                   int       `db:"id_payment_plan"`
	CompanyID            int       `db:"id_company"`
	CompanyName          string    `db:"company_name"`
	Amount               float32   `db:"amount"`
	NumberOfInstallments int       `db:"number_of_installments"`
	Status               string    `db:"status"`
	FirstDueDate         time.Time `db:"first_due_date"`
	LastDueDate          time.Time `db:"last_due_date"`
	UpdatedAt            time.Time `db:"updated_at"`
}

func (r *PaymentPlanRepository) FindAllWithCompany(ctx context.Context) ([]planWithCompanyRow, error) {
	query := `
		SELECT
			PP.id_payment_plan,
			PP.id_company,
			C.name AS company_name,
			PP.amount,
			PP.number_of_installments,
			COALESCE(PP.status, 'pending') AS status,
			PP.first_due_date,
			PP.last_due_date,
			PP.updated_at
		FROM payment_plans PP
		INNER JOIN companies C ON PP.id_company = C.id_company
		WHERE C.deleted_at IS NULL
		ORDER BY C.name ASC, PP.updated_at DESC
	`
	var rows []planWithCompanyRow
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment plans with company: %w", err)
	}
	return rows, nil
}

func (r *PaymentPlanRepository) CountAll(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM payment_plans PP
		INNER JOIN companies C ON PP.id_company = C.id_company
		WHERE C.deleted_at IS NULL
	`
	var total int
	err := r.db.GetContext(ctx, &total, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count payment plans: %w", err)
	}
	return total, nil
}

func (r *PaymentPlanRepository) Insert(ctx context.Context, tx *sqlx.Tx, paymentPlan PaymentPlan) (int, error) {
	query := `
		INSERT INTO payment_plans (
			id_company,
			payments_in_plan,
			amount,
			number_of_installments,
			status,
			first_due_date,
			last_due_date,
			observations
		) VALUES (
			:id_company,
			:payments_in_plan,
			:amount,
			:number_of_installments,
			:status,
			:first_due_date,
			:last_due_date,
			:observations
		)`

	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.NamedExecContext(ctx, query, paymentPlan)
	} else {
		res, err = r.db.NamedExecContext(ctx, query, paymentPlan)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to insert payment plan: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id while inserting payment plan: %w", err)
	}
	return int(id), nil
}

func (r *PaymentPlanRepository) Update(ctx context.Context, id int, paymentPlan PaymentPlan) error {
	query := "UPDATE payment_plans SET observations = ? WHERE id_payment_plan = ?"
	_, err := r.db.ExecContext(ctx, query, paymentPlan.Observations, id)
	if err != nil {
		return fmt.Errorf("failed to update payment plan: %w", err)
	}
	return nil
}

func (r *PaymentPlanRepository) Cancel(ctx context.Context, tx *sqlx.Tx, id int) (int, error) {
	query := "UPDATE payment_plans SET status = 'cancelled' WHERE id_payment_plan = ? AND status = 'pending'"
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to cancel payment plan: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected while cancelling payment plan: %w", err)
	}
	return int(rows), nil
}

func (r *PaymentPlanRepository) Restore(ctx context.Context, tx *sqlx.Tx, id int) (int, error) {
	query := "UPDATE payment_plans SET status = 'pending' WHERE id_payment_plan = ? AND status = 'cancelled'"
	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("failed to restore payment plan: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected while restoring payment plan: %w", err)
	}
	return int(rows), nil
}

func (r *PaymentPlanRepository) UpdateStatus(ctx context.Context, tx *sqlx.Tx, id int, status string) error {
	query := "UPDATE payment_plans SET status = ? WHERE id_payment_plan = ?"
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, status, id)
	} else {
		_, err = r.db.ExecContext(ctx, query, status, id)
	}
	if err != nil {
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

func (r *PaymentPlanRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
