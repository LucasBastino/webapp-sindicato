package auth

import (
	"context"
	"testing"
	"time"

	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
)

type fakeRefreshRepo struct {
	tokens map[string]*authdomain.RefreshToken
	nextID int
}

func newFakeRefreshRepo() *fakeRefreshRepo {
	return &fakeRefreshRepo{tokens: map[string]*authdomain.RefreshToken{}}
}

func (r *fakeRefreshRepo) SaveRefreshToken(_ context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	r.nextID++
	r.tokens[tokenHash] = &authdomain.RefreshToken{
		ID:        r.nextID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
	return nil
}

func (r *fakeRefreshRepo) GetRefreshToken(_ context.Context, tokenHash string) (*authdomain.RefreshToken, error) {
	t, ok := r.tokens[tokenHash]
	if !ok {
		return nil, nil
	}
	copy := *t
	return &copy, nil
}

func (r *fakeRefreshRepo) RevokeRefreshToken(_ context.Context, tokenHash string) error {
	if t, ok := r.tokens[tokenHash]; ok && t.RevokedAt == nil {
		now := time.Now()
		t.RevokedAt = &now
	}
	return nil
}

func (r *fakeRefreshRepo) RevokeAllRefreshTokensForUser(_ context.Context, userID int) error {
	now := time.Now()
	for _, t := range r.tokens {
		if t.UserID == userID && t.RevokedAt == nil {
			t.RevokedAt = &now
		}
	}
	return nil
}

func (r *fakeRefreshRepo) RevokeAllRefreshTokensForUserExcept(_ context.Context, userID int, exceptTokenHash string) error {
	now := time.Now()
	for _, t := range r.tokens {
		if t.UserID == userID && t.RevokedAt == nil && t.TokenHash != exceptTokenHash {
			t.RevokedAt = &now
		}
	}
	return nil
}

func TestRevokeAllRefreshTokensForUser_MarksAll(t *testing.T) {
	fake := newFakeRefreshRepo()
	_ = fake.SaveRefreshToken(context.Background(), 7, "h1", time.Now().Add(time.Hour))
	_ = fake.SaveRefreshToken(context.Background(), 7, "h2", time.Now().Add(time.Hour))
	_ = fake.SaveRefreshToken(context.Background(), 8, "h3", time.Now().Add(time.Hour))
	if err := fake.RevokeAllRefreshTokensForUser(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
	t1, _ := fake.GetRefreshToken(context.Background(), "h1")
	t2, _ := fake.GetRefreshToken(context.Background(), "h2")
	t3, _ := fake.GetRefreshToken(context.Background(), "h3")
	if t1.RevokedAt == nil || t2.RevokedAt == nil {
		t.Fatal("expected user 7 tokens revoked")
	}
	if t3.RevokedAt != nil {
		t.Fatal("did not expect user 8 token revoked")
	}
}

func TestRevokeAllRefreshTokensForUserExcept_KeepsCurrent(t *testing.T) {
	fake := newFakeRefreshRepo()
	_ = fake.SaveRefreshToken(context.Background(), 7, "keep", time.Now().Add(time.Hour))
	_ = fake.SaveRefreshToken(context.Background(), 7, "revoke1", time.Now().Add(time.Hour))
	_ = fake.SaveRefreshToken(context.Background(), 7, "revoke2", time.Now().Add(time.Hour))
	_ = fake.SaveRefreshToken(context.Background(), 8, "other-user", time.Now().Add(time.Hour))
	if err := fake.RevokeAllRefreshTokensForUserExcept(context.Background(), 7, "keep"); err != nil {
		t.Fatal(err)
	}
	keep, _ := fake.GetRefreshToken(context.Background(), "keep")
	r1, _ := fake.GetRefreshToken(context.Background(), "revoke1")
	r2, _ := fake.GetRefreshToken(context.Background(), "revoke2")
	other, _ := fake.GetRefreshToken(context.Background(), "other-user")
	if keep.RevokedAt != nil {
		t.Fatal("expected current token to remain active")
	}
	if r1.RevokedAt == nil || r2.RevokedAt == nil {
		t.Fatal("expected other user 7 tokens revoked")
	}
	if other.RevokedAt != nil {
		t.Fatal("did not expect other user token revoked")
	}
}

func TestRevokeOtherSessionsAfterRotation_OldTokenRevokesNew(t *testing.T) {
	fake := newFakeRefreshRepo()
	rawR1 := "refresh-token-one"
	rawR2 := "refresh-token-two"
	hashR1 := hashToken(rawR1)
	hashR2 := hashToken(rawR2)

	_ = fake.SaveRefreshToken(context.Background(), 1, hashR1, time.Now().Add(time.Hour))
	now := time.Now()
	fake.tokens[hashR1].RevokedAt = &now
	_ = fake.SaveRefreshToken(context.Background(), 1, hashR2, time.Now().Add(time.Hour))

	if err := fake.RevokeAllRefreshTokensForUserExcept(context.Background(), 1, hashR1); err != nil {
		t.Fatal(err)
	}

	neu, _ := fake.GetRefreshToken(context.Background(), hashR2)
	if neu.RevokedAt == nil {
		t.Fatal("using old token hash after rotation should revoke the new session (bug scenario)")
	}
}

func TestRevokeOtherSessionsAfterRotation_NewTokenKeepsNew(t *testing.T) {
	fake := newFakeRefreshRepo()
	rawR1 := "refresh-token-one"
	rawR2 := "refresh-token-two"
	hashR1 := hashToken(rawR1)
	hashR2 := hashToken(rawR2)

	_ = fake.SaveRefreshToken(context.Background(), 1, hashR1, time.Now().Add(time.Hour))
	now := time.Now()
	fake.tokens[hashR1].RevokedAt = &now
	_ = fake.SaveRefreshToken(context.Background(), 1, hashR2, time.Now().Add(time.Hour))

	if err := fake.RevokeAllRefreshTokensForUserExcept(context.Background(), 1, hashR2); err != nil {
		t.Fatal(err)
	}

	neu, _ := fake.GetRefreshToken(context.Background(), hashR2)
	if neu.RevokedAt != nil {
		t.Fatal("using new token hash after rotation should keep the current session active")
	}
}

func TestRefreshRotation_RevokesOldToken(t *testing.T) {
	fake := newFakeRefreshRepo()
	raw := "raw-refresh-token-value"
	hash := hashToken(raw)
	_ = fake.SaveRefreshToken(context.Background(), 1, hash, time.Now().Add(time.Hour))

	// Simulate rotation steps used by AuthService.RefreshAccessToken
	tok, err := fake.GetRefreshToken(context.Background(), hash)
	if err != nil || tok == nil {
		t.Fatal("token missing")
	}
	if err := fake.RevokeRefreshToken(context.Background(), hash); err != nil {
		t.Fatal(err)
	}
	newRaw := "new-raw-refresh-token"
	newHash := hashToken(newRaw)
	if err := fake.SaveRefreshToken(context.Background(), 1, newHash, time.Now().Add(time.Hour)); err != nil {
		t.Fatal(err)
	}

	old, _ := fake.GetRefreshToken(context.Background(), hash)
	neu, _ := fake.GetRefreshToken(context.Background(), newHash)
	if old.RevokedAt == nil {
		t.Fatal("old refresh should be revoked after rotation")
	}
	if neu.RevokedAt != nil {
		t.Fatal("new refresh should be active")
	}

	// reuse of old revoked token should trigger full revoke
	if old.RevokedAt != nil {
		_ = fake.RevokeAllRefreshTokensForUser(context.Background(), old.UserID)
	}
	neu2, _ := fake.GetRefreshToken(context.Background(), newHash)
	if neu2.RevokedAt == nil {
		t.Fatal("reuse of rotated token should revoke all user sessions")
	}
}
