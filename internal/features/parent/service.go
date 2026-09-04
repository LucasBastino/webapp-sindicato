package parent

import (
	"context"
	"errors"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
)

type ParentService struct {
	repo        *ParentRepository
	idempotency *idempotency.IdempotencyService
}

func NewParentService(repo *ParentRepository, idempotencyService *idempotency.IdempotencyService) *ParentService {
	return &ParentService{
		repo:        repo,
		idempotency: idempotencyService,
	}
}


func (s *ParentService) Get(ctx context.Context, id int) (*Parent, error) {
	parent, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if parent == nil{
		return nil, apperrors.NewNotFoundError(errors.New("parent not found"), "")
	}

	return parent, nil
}

func (s *ParentService) List(ctx context.Context, memberID int) ([]Parent, error) {
	parents, err := s.repo.FindAll(ctx, memberID)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return parents, nil
}

func (s *ParentService) Count(ctx context.Context, memberID int) (int, error) {
	count, err := s.repo.Count(ctx, memberID)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	return count, nil
}


func (s *ParentService) Create(ctx context.Context, parent Parent, idempotencyKey string) (int, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx creating parent: %w", err), "")
	}
	defer tx.Rollback()

	id, err := s.repo.Insert(ctx, tx, parent)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}
	if id == 0 {
		return 0, apperrors.NewBusinessError(errors.New("cannot create parent: member is inactive or does not exist"), "El afiliado no se encuentra activo.")
	}

	if err := s.idempotency.UpdateResource(ctx, tx, idempotencyKey, "parent", id); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to commit parent create: %w", err), "")
	}

	return id, nil
}


func (s *ParentService) Update(ctx context.Context, id int, parent Parent) error {
	err := s.repo.Update(ctx, id, parent)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *ParentService) HardDelete(ctx context.Context, id int) error {
	err := s.repo.HardDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}




// func (s *ParentService) ValidateInDBForCreate(ctx context.Context, memberID int, inputs Request) map[string]string {
// 	return s.repo.ValidateInDB(ctx, 0, memberID, inputs)
// }

// func (s *ParentService) ValidateInDBforUpdate(ctx context.Context, parentID, memberID int, inputs Request) map[string]string {
// 	return s.repo.ValidateInDB(ctx, parentID, memberID, inputs)
// }
