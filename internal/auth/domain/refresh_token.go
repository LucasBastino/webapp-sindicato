package authdomain

import "time"

type RefreshToken struct {
	ID        int        `db:"id"`
	UserID    int        `db:"id_user"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	RevokedAt *time.Time `db:"revoked_at"`
	CreatedAt time.Time  `db:"created_at"`
}