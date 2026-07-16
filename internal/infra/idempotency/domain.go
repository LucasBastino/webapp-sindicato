package idempotency

import "time"

type Record struct {
	ID             int       `db:"id"`

	IdempotencyKey string    `db:"idempotency_key"`
	RequestHash    string    `db:"request_hash"`

	ResourceType   *string   `db:"resource_type"`
	ResourceID     *int      `db:"resource_id"`
    
	CreatedAt      time.Time `db:"created_at"`
	ExpiresAt      time.Time `db:"expires_at"`
}