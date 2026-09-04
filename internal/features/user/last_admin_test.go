package user

import (
	"context"
	"errors"
	"testing"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
)

func TestEnsureNotLastAdminSkipsNonAdmin(t *testing.T) {
	s := &UserService{}
	err := s.ensureNotLastAdmin(context.Background(), &User{ID: 1, Admin: false})
	if err != nil {
		t.Fatalf("expected nil for non-admin, got %v", err)
	}
}

func TestEnsureNotLastAdminSkipsNilUser(t *testing.T) {
	s := &UserService{}
	err := s.ensureNotLastAdmin(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected nil for nil user, got %v", err)
	}
}

func TestCannotRemoveLastAdminSentinel(t *testing.T) {
	err := apperrors.NewBusinessError(
		apperrors.ErrCannotRemoveLastAdmin,
		"No se puede eliminar ni degradar al último administrador.",
	)
	if !errors.Is(err, apperrors.ErrCannotRemoveLastAdmin) {
		t.Fatal("BusinessError should unwrap to ErrCannotRemoveLastAdmin")
	}
}
