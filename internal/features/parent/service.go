package parent

import (
	"context"
	"errors"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
)

type ParentService struct {
	repo *ParentRepository
}

func NewParentService(repo *ParentRepository) *ParentService {
	return &ParentService{repo: repo}
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


func (s *ParentService) Create(ctx context.Context, parent Parent) (int, error) {
	id, err := s.repo.Insert(ctx, parent)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
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
