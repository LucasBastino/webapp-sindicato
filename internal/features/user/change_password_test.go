package user

import (
	"context"
	"testing"

	"github.com/LucasBastino/webapp-sindicato/internal/security/password"
)

func TestChangePasswordSelfRevokesOtherSessionsOnly(t *testing.T) {
	hasher := password.NewBcryptHasher()
	hash, err := hasher.Generate("oldpass")
	if err != nil {
		t.Fatal(err)
	}

	revoker := &recordingRevoker{}
	svc := &UserService{
		repo: &stubUserRepo{
			user: &User{ID: 1, PasswordHash: string(hash)},
		},
		passwordHasher: hasher,
		sessionRevoker: revoker,
	}

	err = svc.ChangePassword(context.Background(), 1, 1, "current-refresh-token", passwordRequest{
		CurrentPassword: "oldpass",
		Password:        "newpass",
		ConfirmPassword: "newpass",
	}, true)
	if err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}
	if len(revoker.allCalled) != 0 {
		t.Fatalf("did not expect revoke all, got %#v", revoker.allCalled)
	}
	if len(revoker.otherCalled) != 1 {
		t.Fatalf("expected one revoke-other call, got %#v", revoker.otherCalled)
	}
	if revoker.otherCalled[0].userID != 1 || revoker.otherCalled[0].token != "current-refresh-token" {
		t.Fatalf("unexpected revoke-other payload: %#v", revoker.otherCalled[0])
	}
}

func TestChangePasswordOtherUserRevokesAllSessions(t *testing.T) {
	hasher := password.NewBcryptHasher()
	hash, err := hasher.Generate("oldpass")
	if err != nil {
		t.Fatal(err)
	}

	revoker := &recordingRevoker{}
	svc := &UserService{
		repo: &stubUserRepo{
			user: &User{ID: 2, PasswordHash: string(hash)},
		},
		passwordHasher: hasher,
		sessionRevoker: revoker,
	}

	err = svc.ChangePassword(context.Background(), 2, 1, "", passwordRequest{
		Password:        "newpass",
		ConfirmPassword: "newpass",
	}, false)
	if err != nil {
		t.Fatalf("ChangePassword: %v", err)
	}
	if len(revoker.allCalled) != 1 || revoker.allCalled[0] != 2 {
		t.Fatalf("expected revoke all for user 2, got %#v", revoker.allCalled)
	}
	if len(revoker.otherCalled) != 0 {
		t.Fatalf("did not expect revoke-other, got %#v", revoker.otherCalled)
	}
}
