package company

import (
	"context"
	"errors"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/common/page"
	"github.com/LucasBastino/webapp-sindicato/internal/features/payment"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
)

type CompanyService struct{
	repo *CompanyRepository

	paymentService *payment.PaymentService
	idempotency    *idempotency.IdempotencyService

	logger logger.Logger
}

func NewCompanyService(repo *CompanyRepository, paymentService *payment.PaymentService, idempotencyService *idempotency.IdempotencyService) *CompanyService{
	return &CompanyService{
		repo:           repo,
		paymentService: paymentService,
		idempotency:    idempotencyService,
	}
}

func (s *CompanyService) GetName(ctx context.Context, id int) (string, error) {
	company, err := s.Get(ctx, id)
	if err != nil {
		return "", err
	}
	return company.Name, nil
}

func (s *CompanyService) GetNavDetails(ctx context.Context, id int) (name, number, address, phone string, err error) {
	company, err := s.Get(ctx, id)
	if err != nil {
		return "", "", "", "", err
	}
	number = ""
	if company.CompanyNumber != nil {
		number = *company.CompanyNumber
	}
	return company.Name, number, company.Address, company.Phone, nil
}

func (s *CompanyService) ListActiveIDs(ctx context.Context) ([]int, error) {
	ids, err := s.repo.ListActiveIDs(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	return ids, nil
}

func (s *CompanyService) Get(ctx context.Context, id int) (*Company, error){
	company, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if company == nil{
		return nil, apperrors.NewNotFoundError(errors.New("company not found"), "")
	}

	return company, nil
}

func (s *CompanyService) List(ctx context.Context, filters companyFilters, offset int) ([]Company, error){
	companies, err := s.repo.Search(ctx, filters, offset)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return companies, nil
}

func (s *CompanyService) ListForSelect(ctx context.Context, searchKey string) ([]Company, error){
	companies, err := s.repo.SearchForSelect(ctx, searchKey)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return companies, nil
}

func (s *CompanyService) Count(ctx context.Context, filters companyFilters) (int, error){
	count, err := s.repo.Count(ctx, filters)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	return count, nil
}

func (s *CompanyService) CountActive(ctx context.Context) (int, error) {
	return s.Count(ctx, companyFilters{
		statuses: page.StatusFilters{ShowActive: true},
	})
}

func (s *CompanyService) FindRecent(ctx context.Context, limit int) ([]RecentCompany, error) {
	companies, err := s.repo.FindRecent(ctx, limit)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	return companies, nil
}


func (s *CompanyService) Create(ctx context.Context, company Company, idempotencyKey string) (int, error){
	tx, err := s.repo.BeginTx(ctx)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to iniciate transaction while creating company: %w", err), "")
	}

	defer tx.Rollback()

	defer func(){
		p := recover()
		if p != nil{
			s.logger.Error("panic creating company", "panic", p)
			panic(p)
		}
	}()
	
	id, err := s.repo.Insert(ctx, tx, company)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	// el month desde hoy y el año el de hoy
	err = s.paymentService.CreateRemainingPaymentsByID(ctx, tx, int(id))
	if err!=nil{
		return 0, err 
	}

	if err := s.idempotency.UpdateResource(ctx, tx, idempotencyKey, "company", id); err != nil {
		return 0, err
	}
	
	if err := tx.Commit(); err!=nil{
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to commit transaction while creating company: %w", err), "")
	}
	return id, nil
}

func (s *CompanyService) Update(ctx context.Context, id int, company Company) error{
	err := s.repo.Update(ctx, id, company)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *CompanyService) SoftDelete(ctx context.Context, id int) error {
	rows, err := s.repo.SoftDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if rows == 0{
		return apperrors.NewBusinessError(fmt.Errorf("failed to soft delete company: the entity is already inactive or doesn't exist"), "No se pudo completar la operación.")
	}

	return nil
}

func (s *CompanyService) Restore(ctx context.Context, id int) error{
	tx, err := s.repo.BeginTx(ctx)
	if err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to iniciate transaction while restoring company: %w", err), "")
	}

	defer tx.Rollback()

	defer func(){
		p := recover()
		if p != nil{
			s.logger.Error("panic creating company", "panic", p)
			panic(p)
		}
	}()

	rows, err := s.repo.Restore(ctx, tx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0{
		return apperrors.NewBusinessError(errors.New("failed to restore company: the entity is already active or doesn't exist"), "No se pudo completar la operación.")
	}

	// el month desde hoy y el año el de hoy
	err = s.paymentService.CreateRemainingPaymentsByID(ctx, tx, int(id))
	if err!=nil{
		return err
	}

	if err := tx.Commit(); err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to commit transaction while creating company: %w", err), "")
	}
	return nil
}

func (s *CompanyService) HardDelete(ctx context.Context, id int) error {
	err := s.repo.HardDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	
	return nil
}




// func (s *CompanyService) ListMembers(ctx context.Context, id, limit, offset int, searchKey string, includeInactive, includeDeleted bool) ([]member.Member, error){
// 	return s.memberService.ListByCompany(ctx, id, limit, offset, searchKey, includeInactive, includeDeleted)
// }


// func (s *CompanyService) SoftDeleteMember(ctx context.Context, memberID int) error {
// 	return s.memberService.SoftDelete(ctx, memberID)
// }

// func (s *CompanyService) RestoreMember(ctx context.Context, memberID int) error{
// 	return s.memberService.Restore(ctx, memberID)
// }

// func (s *CompanyService) CountMembers(ctx context.Context, id, limit, offset int, searchKey string, includeInactive, includeDeleted bool) (int, error){
// 	return s.memberService.CountByCompany(ctx, id, searchKey, includeInactive, includeDeleted)
// }

func (s *CompanyService) ListPaymentYears(ctx context.Context, id int) ([]int, error){
	return s.paymentService.ListPaymentYears(ctx, id)
}







// func (s *CompanyService) ValidateInDBForCreate(ctx context.Context, inputs Request) map[string]string{
// 	// le paso 0 como id porque todavia no tiene id, lo estoy creando
// 	return s.repo.ValidateCompanyInDB(ctx, 0, inputs)
// }

// func (s *CompanyService) ValidateInDBForUpdate(ctx context.Context, id int, inputs Request) map[string]string{
// 	return s.repo.ValidateCompanyInDB(ctx, id, inputs)
// }

