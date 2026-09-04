package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository{
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
    query := `INSERT INTO refresh_tokens (id_user, token_hash, expires_at) VALUES (?, ?, ?)`
    _, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt)
    if err != nil {
        return fmt.Errorf("failed to save refresh token: %w", err)
    }
    return nil
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, tokenHash string) (*authdomain.RefreshToken, error) {
    query := `SELECT * FROM refresh_tokens WHERE token_hash = ?`
    var refreshToken authdomain.RefreshToken
    err := r.db.GetContext(ctx, &refreshToken, query, tokenHash)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get refresh token: %w", err)
    }
    return &refreshToken, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
    query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE token_hash = ? AND revoked_at IS NULL`
    _, err := r.db.ExecContext(ctx, query, tokenHash)
    if err != nil {
        return fmt.Errorf("failed to revoke refresh token: %w", err)
    }
    return nil
}

func (r *AuthRepository) RevokeAllRefreshTokensForUser(ctx context.Context, userID int) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id_user = ? AND revoked_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens for user %d: %w", userID, err)
	}
	return nil
}

func (r *AuthRepository) RevokeAllRefreshTokensForUserExcept(ctx context.Context, userID int, exceptTokenHash string) error {
	query := `UPDATE refresh_tokens SET revoked_at = NOW() WHERE id_user = ? AND revoked_at IS NULL AND token_hash != ?`
	_, err := r.db.ExecContext(ctx, query, userID, exceptTokenHash)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens for user %d except current: %w", userID, err)
	}
	return nil
}
