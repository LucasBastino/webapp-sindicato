package user

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
)

type recordingRevoker struct {
	allCalled   []int
	otherCalled []otherRevokeCall
}

type otherRevokeCall struct {
	userID int
	token  string
}

func (r *recordingRevoker) RevokeAllSessionsForUser(_ context.Context, userID int) error {
	r.allCalled = append(r.allCalled, userID)
	return nil
}

func (r *recordingRevoker) RevokeOtherSessionsForUser(_ context.Context, userID int, currentRefreshToken string) error {
	r.otherCalled = append(r.otherCalled, otherRevokeCall{userID: userID, token: currentRefreshToken})
	return nil
}

type stubUserRepo struct {
	user       *User
	adminCount int
}

func (s *stubUserRepo) FindByID(_ context.Context, id int) (*User, error) {
	if s.user == nil {
		return nil, nil
	}
	copy := *s.user
	copy.ID = id
	return &copy, nil
}

func (s *stubUserRepo) FindByUsername(context.Context, string) (*User, error) { return nil, nil }
func (s *stubUserRepo) FindAll(context.Context) ([]User, error)               { return nil, nil }
func (s *stubUserRepo) CountAdmins(context.Context) (int, error)              { return s.adminCount, nil }
func (s *stubUserRepo) UpdatePermissions(context.Context, int, bool, map[string]any) error {
	return nil
}
func (s *stubUserRepo) UpdatePassword(context.Context, int, string) (int, error) { return 1, nil }
func (s *stubUserRepo) HardDelete(context.Context, int) error                    { return nil }
func (s *stubUserRepo) BeginTx(context.Context) (*sqlx.Tx, error)                { return nil, nil }
func (s *stubUserRepo) Insert(context.Context, *sqlx.Tx, User) (int, error)      { return 0, nil }

func TestUpdatePermissionsRevokesSessions(t *testing.T) {
	revoker := &recordingRevoker{}
	svc := &UserService{
		repo: &stubUserRepo{
			user:       &User{ID: 2, Admin: false},
			adminCount: 1,
		},
		sessionRevoker: revoker,
	}

	err := svc.UpdatePermissions(context.Background(), 2, 1, false, map[string]any{
		"member":  "viewer",
		"company": "viewer",
	})
	if err != nil {
		t.Fatalf("UpdatePermissions: %v", err)
	}
	if len(revoker.allCalled) != 1 || revoker.allCalled[0] != 2 {
		t.Fatalf("expected revoke of user 2, got %#v", revoker.allCalled)
	}
}

func TestUpdatePermissionsDoesNotRevokeOwnEdit(t *testing.T) {
	revoker := &recordingRevoker{}
	svc := &UserService{sessionRevoker: revoker}
	err := svc.UpdatePermissions(context.Background(), 1, 1, false, nil)
	if err == nil {
		t.Fatal("expected cannot edit own permissions")
	}
	if len(revoker.allCalled) != 0 {
		t.Fatalf("did not expect revoke, got %#v", revoker.allCalled)
	}
}

func TestHardDeleteRevokesSessions(t *testing.T) {
	revoker := &recordingRevoker{}
	svc := &UserService{
		repo: &stubUserRepo{
			user:       &User{ID: 2, Admin: false},
			adminCount: 1,
		},
		sessionRevoker: revoker,
	}
	if err := svc.HardDelete(context.Background(), 2, 1); err != nil {
		t.Fatalf("HardDelete: %v", err)
	}
	if len(revoker.allCalled) != 1 || revoker.allCalled[0] != 2 {
		t.Fatalf("expected revoke of user 2, got %#v", revoker.allCalled)
	}
}
