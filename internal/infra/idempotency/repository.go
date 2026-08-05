package idempotency

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type IdempotencyRepository struct {
	db *sqlx.DB
}

func NewIdempotencyRepository(db *sqlx.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (r *IdempotencyRepository) FindByKey(ctx context.Context, key string) (*Record, error) {

	query := `
	SELECT *
	FROM idempotency_keys
	WHERE idempotency_key = ?
	`

	var record Record

	err := r.db.GetContext(ctx, &record, query, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch idempotency key: %w", err)
	}

	return &record, nil
}

func (r *IdempotencyRepository) Create(ctx context.Context,	key string, requestHash string,	expiresAt time.Time) error {

	query := `
	INSERT INTO idempotency_keys (
		idempotency_key,
		request_hash,
		expires_at
	)
	VALUES (?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx,	query, key,	requestHash, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to create idempotency key: %w", err)
	}

	return nil
}

func (r *IdempotencyRepository) UpdateResource(ctx context.Context, tx *sqlx.Tx, key string, resourceType string, resourceID int) error {
	query := `
	UPDATE idempotency_keys
	SET
		resource_type = ?,
		resource_id = ?
	WHERE
		idempotency_key = ?
	`

	_, err := tx.ExecContext(ctx, query, resourceType, resourceID, key)
	if err != nil {
		return fmt.Errorf("failed to update idempotency resource: %w", err)
	}

	return nil
}

func (r *IdempotencyRepository) DeleteExpired(ctx context.Context) error {
	query := "DELETE FROM idempotency_keys WHERE expires_at < NOW()"
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired idempotency keys: %w", err)
	}

	return nil
}