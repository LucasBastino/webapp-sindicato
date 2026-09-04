package member

import (
	"context"
	"errors"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/common/page"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
)

type MemberService struct {
	repo        *MemberRepository
	idempotency *idempotency.IdempotencyService
}

func NewMemberService(repo *MemberRepository, idempotencyService *idempotency.IdempotencyService) *MemberService {
	return &MemberService{
		repo:        repo,
		idempotency: idempotencyService,
	}
}

func (s *MemberService) Get(ctx context.Context, id int) (*Member, error) {
	member, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if member == nil{
		return nil, apperrors.NewNotFoundError(errors.New("member not found"), "")
	}

	return member, nil
}

func (s *MemberService) List(ctx context.Context, filters memberFilters, offset int) ([]Member, error) {
	members, err := s.repo.Search(ctx, filters, offset)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return members, nil
}

func (s *MemberService) GetElectoralList(ctx context.Context) ([]Member, error) {
	members, err := s.repo.GetElectoralMemberList(ctx)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return members, nil
}

func (s *MemberService) Count(ctx context.Context, filters memberFilters) (int, error) {
	count, err := s.repo.Count(ctx, filters)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	return count, nil
}

func (s *MemberService) CountActive(ctx context.Context) (int, error) {
	return s.Count(ctx, memberFilters{
		statuses: page.StatusFilters{ShowActive: true},
	})
}

func (s *MemberService) FindRecent(ctx context.Context, limit int) ([]Member, error) {
	members, err := s.repo.FindRecent(ctx, limit)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	return members, nil
}

func (s *MemberService) Create(ctx context.Context, member Member, idempotencyKey string) (int, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx creating member: %w", err), "")
	}
	defer tx.Rollback()

	id, err := s.repo.Insert(ctx, tx, member)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}
	if id == 0 {
		return 0, apperrors.NewBusinessError(errors.New("cannot create member: company is inactive or does not exist"), "La empresa no se encuentra activa.")
	}

	if err := s.idempotency.UpdateResource(ctx, tx, idempotencyKey, "member", id); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to commit member create: %w", err), "")
	}

	return id, nil
}

func (s *MemberService) Update(ctx context.Context, id int, member Member) error {
	err := s.repo.Update(ctx, id, member)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *MemberService) SoftDelete(ctx context.Context, id int) error {
	rows, err := s.repo.SoftDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0 {
		return apperrors.NewBusinessError(errors.New("failed to softdelete member: the entity is already softdeleted or doesn't exist"), "No se pudo completar la operación.")
	}
	return nil
}

func (s *MemberService) Restore(ctx context.Context, id int) error {
	rows, err := s.repo.Restore(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0{
		return apperrors.NewBusinessError(errors.New("failed to restore member: the entity is already active or doesn't exist"), "No se pudo completar la operación.")
	}
	return nil
}

// borra definitivamente el afiliado y sus parientes de la base de datos
// no hace falta chequear rows porque es una funcion idempotente
func (s *MemberService) HardDelete(ctx context.Context, id int) error {
	err := s.repo.HardDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	return nil
}
