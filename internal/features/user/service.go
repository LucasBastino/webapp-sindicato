package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/security/password"
)
type UserService struct{
	repo *UserRepository
	
	passwordHasher password.Hasher
}

func NewUserService(repo *UserRepository, passwordHasher password.Hasher) *UserService{
	return &UserService{
		repo: repo,
		passwordHasher: passwordHasher,
	}
}


func (s *UserService) Create(ctx context.Context, user User) (int, error) {

	if user.Admin{
		for _, role := range user.ResourceRoles{
			if role!= "editor"{
				return 0, apperrors.NewBusinessError(apperrors.ErrInvalidPermissions, "")
			}
		}
	}

	id, err := s.repo.Insert(ctx, user)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
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


func (s *UserService) UpdatePermissions(ctx context.Context, id int, admin bool, resourceRoles map[string]any) error {
	if admin{
		for _, role := range resourceRoles{
			if role!= "editor"{
				return apperrors.NewBusinessError(apperrors.ErrInvalidPermissions, "")
			}
		}
	}
	err := s.repo.UpdatePermissions(ctx, id, admin, resourceRoles)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *UserService) ChangePassword(ctx context.Context, id int, req passwordRequest) error {
	user, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if user == nil{
		return apperrors.NewNotFoundError(errors.New("user not found"), "")
	}

	err = s.passwordHasher.Compare([]byte(user.PasswordHash), req.CurrentPassword)
	if err!=nil{
		return apperrors.NewBusinessError(apperrors.ErrInvalidCurrentPassword, "")
	}

	if req.CurrentPassword == req.Password {
		return apperrors.NewBusinessError(apperrors.ErrInvalidNewPassword, "")
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

func (s *UserService) HardDelete(ctx context.Context, id int) error {
	err := s.repo.HardDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}
