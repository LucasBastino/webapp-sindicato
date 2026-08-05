package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/app-sindicato/internal/security/password"
)
type UserService struct{
	repo *UserRepository
	
	passwordHasher password.Hasher
	idempotency    *idempotency.IdempotencyService
}

func NewUserService(repo *UserRepository, passwordHasher password.Hasher, idempotencyService *idempotency.IdempotencyService) *UserService{
	return &UserService{
		repo:           repo,
		passwordHasher: passwordHasher,
		idempotency:    idempotencyService,
	}
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


func (s *UserService) UpdatePermissions(ctx context.Context, id int, actorID int, admin bool, resourceRoles map[string]any) error {
	if id == actorID {
		return apperrors.NewBusinessError(apperrors.ErrCannotEditOwnPermissions, "No podés editar tus propios permisos.")
	}

	if admin {
		resourceRoles = map[string]any{
			"member":  "editor",
			"company": "editor",
		}
	}

	err := s.repo.UpdatePermissions(ctx, id, admin, resourceRoles)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, id int, req passwordRequest, requireCurrent bool) error {
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

	return nil
}

func (s *UserService) HardDelete(ctx context.Context, id int, actorID int) error {
	if id == actorID {
		return apperrors.NewBusinessError(apperrors.ErrCannotDeleteSelf, "No podés eliminarte a vos mismo.")
	}

	err := s.repo.HardDelete(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}
