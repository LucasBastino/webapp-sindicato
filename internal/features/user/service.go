package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/webapp-sindicato/internal/security/password"
	"github.com/jmoiron/sqlx"
)
type UserService struct{
	repo userRepo
	
	passwordHasher password.Hasher
	idempotency    *idempotency.IdempotencyService
	sessionRevoker SessionRevoker
}

type userRepo interface {
	FindByID(ctx context.Context, id int) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	CountAdmins(ctx context.Context) (int, error)
	UpdatePermissions(ctx context.Context, id int, admin bool, resourceRoles map[string]any) error
	UpdatePassword(ctx context.Context, id int, hash string) (int, error)
	HardDelete(ctx context.Context, id int) error
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	Insert(ctx context.Context, tx *sqlx.Tx, user User) (int, error)
}

// SessionRevoker invalidates auth sessions (refresh tokens) for a user.
type SessionRevoker interface {
	RevokeAllSessionsForUser(ctx context.Context, userID int) error
	RevokeOtherSessionsForUser(ctx context.Context, userID int, currentRefreshToken string) error
}

func NewUserService(repo *UserRepository, passwordHasher password.Hasher, idempotencyService *idempotency.IdempotencyService) *UserService{
	return &UserService{
		repo:           repo,
		passwordHasher: passwordHasher,
		idempotency:    idempotencyService,
	}
}

func (s *UserService) SetSessionRevoker(revoker SessionRevoker) {
	s.sessionRevoker = revoker
}


func (s *UserService) Create(ctx context.Context, user User, idempotencyKey string) (int, error) {
	if user.Admin {
		user.ResourceRoles = map[string]any{
			"member":  "editor",
			"company": "editor",
		}
	}

	rolesJSON, err := json.Marshal(user.ResourceRoles)
	if err != nil {
		return 0, apperrors.NewInternalError(fmt.Errorf("failed to marshal resource roles: %w", err), "")
	}
	user.ResourceRolesJSON = rolesJSON

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx creating user: %w", err), "")
	}
	defer tx.Rollback()

	id, err := s.repo.Insert(ctx, tx, user)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}

	if err := s.idempotency.UpdateResource(ctx, tx, idempotencyKey, "user", id); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to commit user create: %w", err), "")
	}

	return id, nil
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.repo.FindByUsername(ctx, username)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return user, nil
}

func (s *UserService) Get(ctx context.Context, id int) (*User, error){
	user, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if user == nil{
		return nil, apperrors.NewNotFoundError(errors.New("user not found"), "")
	}

	return user, nil
}

func (s *UserService) List(ctx context.Context) ([]User, error) {
    users, err := s.repo.FindAll(ctx)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return users, nil
}


func (s *UserService) ensureNotLastAdmin(ctx context.Context, user *User) error {
	if user == nil || !user.Admin {
		return nil
	}
	count, err := s.repo.CountAdmins(ctx)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if count <= 1 {
		return apperrors.NewBusinessError(
			apperrors.ErrCannotRemoveLastAdmin,
			"No se puede eliminar ni degradar al último administrador.",
		)
	}
	return nil
}

func (s *UserService) UpdatePermissions(ctx context.Context, id int, actorID int, admin bool, resourceRoles map[string]any) error {
	if id == actorID {
		return apperrors.NewBusinessError(apperrors.ErrCannotEditOwnPermissions, "No podés editar tus propios permisos.")
	}

	user, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if !admin {
		if err := s.ensureNotLastAdmin(ctx, user); err != nil {
			return err
		}
	}

	if admin {
		resourceRoles = map[string]any{
			"member":  "editor",
			"company": "editor",
		}
	}

	err = s.repo.UpdatePermissions(ctx, id, admin, resourceRoles)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return s.revokeUserSessions(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, id int, actorID int, currentRefreshToken string, req passwordRequest, requireCurrent bool) error {
	user, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if user == nil{
		return apperrors.NewNotFoundError(errors.New("user not found"), "")
	}

	if requireCurrent {
		err = s.passwordHasher.Compare([]byte(user.PasswordHash), req.CurrentPassword)
		if err!=nil{
			return apperrors.NewBusinessError(apperrors.ErrInvalidCurrentPassword, "")
		}

		if req.CurrentPassword == req.Password {
			return apperrors.NewBusinessError(apperrors.ErrInvalidNewPassword, "")
		}
	} else {
		err = s.passwordHasher.Compare([]byte(user.PasswordHash), req.Password)
		if err == nil {
			return apperrors.NewBusinessError(apperrors.ErrInvalidNewPassword, "")
		}
	}

    hash, err := s.passwordHasher.Generate(req.Password)
    if err != nil {
        return apperrors.NewInternalError(fmt.Errorf("failed to generate password hash: %w", err), "")
    }

    rows, err := s.repo.UpdatePassword(ctx, id, string(hash))
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if rows == 0{
		return apperrors.NewNotFoundError(fmt.Errorf("failed to update user password: user doesn't exist"), "El usuario que quieres modificar no existe.")
	}

	if actorID == id {
		return s.revokeOtherUserSessions(ctx, id, currentRefreshToken)
	}
	return s.revokeUserSessions(ctx, id)
}

func (s *UserService) HardDelete(ctx context.Context, id int, actorID int) error {
	if id == actorID {
		return apperrors.NewBusinessError(apperrors.ErrCannotDeleteSelf, "No podés eliminarte a vos mismo.")
	}

	user, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := s.ensureNotLastAdmin(ctx, user); err != nil {
		return err
	}

	err = s.repo.HardDelete(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return s.revokeUserSessions(ctx, id)
}

func (s *UserService) revokeUserSessions(ctx context.Context, userID int) error {
	if s.sessionRevoker == nil {
		return nil
	}
	return s.sessionRevoker.RevokeAllSessionsForUser(ctx, userID)
}

func (s *UserService) revokeOtherUserSessions(ctx context.Context, userID int, currentRefreshToken string) error {
	if s.sessionRevoker == nil {
		return nil
	}
	return s.sessionRevoker.RevokeOtherSessionsForUser(ctx, userID, currentRefreshToken)
}
