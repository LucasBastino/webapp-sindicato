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

func NewInstallmentRepository(db *sqlx.DB) *InstallmentRepository {
	return &InstallmentRepository{db: db}
}

func (r *InstallmentRepository) FindByID(ctx context.Context, id int) (*Installment, error) {
	query := `
		SELECT
			id_installment,
			id_payment_plan,
			installment_number,
			amount,
			due_date,
			paid_at,
			COALESCE(observations, '') AS observations,
			updated_at
		FROM installments
		WHERE id_installment = ?
	`
	var installment Installment
	err := r.db.GetContext(ctx, &installment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch installment: %w", err)
	}
	return &installment, nil
}

func (r *InstallmentRepository) FindAll(ctx context.Context, tx *sqlx.Tx, paymentPlanID int) ([]Installment, error) {
	query := `
		SELECT
			id_installment,
			id_payment_plan,
			installment_number,
			amount,
			due_date,
			paid_at,
			COALESCE(observations, '') AS observations,
			updated_at
		FROM installments
		WHERE id_payment_plan = ?
		ORDER BY installment_number ASC
	`
	var installments []Installment
	var err error
	if tx != nil {
		err = tx.SelectContext(ctx, &installments, query, paymentPlanID)
	} else {
		err = r.db.SelectContext(ctx, &installments, query, paymentPlanID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch installments: %w", err)
	}
	return installments, nil
}

func (r *InstallmentRepository) BulkInsert(ctx context.Context, tx *sqlx.Tx, installments []Installment) (int, error) {
	if len(installments) == 0 {
		return 0, nil
	}
	query := `
	INSERT INTO installments (
		id_payment_plan,
		installment_number,
		amount,
		due_date,
		observations
	) VALUES
	`
	placeholders := make([]string, 0, len(installments))
	args := make([]any, 0, len(installments)*5)

	for _, i := range installments {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?)")
		args = append(args, i.PaymentPlanID, i.InstallmentNumber, i.Amount, i.DueDate, i.Observations)
	}

	query += strings.Join(placeholders, ",")

	var res sql.Result
	var err error
	if tx != nil {
		res, err = tx.ExecContext(ctx, query, args...)
	} else {
		res, err = r.db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to bulk insert installments: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected while bulk inserting installments: %w", err)
	}
	return int(rows), nil
}

func (r *InstallmentRepository) Update(ctx context.Context, tx *sqlx.Tx, id int, installment Installment) error {
	installment.ID = id
	query := `
	UPDATE installments
	SET
		paid_at = :paid_at,
		observations = :observations
	WHERE
		id_installment = :id_installment
	`
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, installment)
	} else {
		_, err = r.db.NamedExecContext(ctx, query, installment)
	}
	if err != nil {
		return fmt.Errorf("failed to update installment: %w", err)
	}
	return nil
}
