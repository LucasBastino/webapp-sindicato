package paymentplan

import (
	"context"
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/features/installment"
	"github.com/LucasBastino/app-sindicato/internal/features/payment"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/jmoiron/sqlx"
)

type PaymentPlanService struct {
	repo *PaymentPlanRepository
	
	paymentRepo *payment.PaymentRepository

	installmentRepo *installment.InstallmentRepository

	logger logger.Logger
}

func NewPaymentPlanService(repo *PaymentPlanRepository, paymentRepo *payment.PaymentRepository, installmentRepo *installment.InstallmentRepository, logger logger.Logger) *PaymentPlanService {
	return &PaymentPlanService{
		repo: repo,
		paymentRepo: paymentRepo,
		installmentRepo: installmentRepo,
		logger: logger,
	}
}

func (s *PaymentPlanService) Get(ctx context.Context, id int) (*PaymentPlan, []installment.Installment, error) {
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, nil, apperrors.NewDatabaseError(err, "")
	}

	if paymentPlan == nil{
		return nil, nil, apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}

	installments, err := s.installmentRepo.FindAll(ctx, nil, id)
	if err!=nil{
		return nil, nil, apperrors.NewDatabaseError(err, "")
	}

	return paymentPlan, installments, err
}

func (s *PaymentPlanService) List(ctx context.Context, companyID int) ([]PaymentPlan, error) {
	paymentPlans, err := s.repo.FindAll(ctx, companyID)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return paymentPlans, nil
}

func (s *PaymentPlanService) Create(ctx context.Context, paymentPlan PaymentPlan) (int, error) {
	id, err := s.repo.Insert(ctx, paymentPlan)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	return id, nil
}

func (s *PaymentPlanService) Update(ctx context.Context, id int, paymentPlan PaymentPlan) error {
	err := s.repo.Update(ctx, id, paymentPlan)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *PaymentPlanService) Cancel(ctx context.Context, id int) error {
	tx, err := s.repo.BeginTx(ctx)
	if err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to iniciate transaction while cancelling payment plan: %w", err), "")
	}

	defer tx.Rollback()

	defer func(){
		p := recover()
		if p != nil{
			s.logger.Error("panic cancelling payment plan", "panic", p)
			panic(p)
		}
	}()
	
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if paymentPlan == nil{
		return apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}


	rows, err := s.repo.Cancel(ctx, tx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0{
		return apperrors.NewBusinessError(errors.New("failed to cancel payment plan: is already cancelled or doesn't exist"), "El plan de pago ya está cancelado o no existe." )
	}

	rows, err = s.paymentRepo.BulkUpdateStatus(ctx, tx, paymentPlan.PaymentsInPlan, "overdue")
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	err = tx.Commit()
	if err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to commit transaction while cancelling payment plan: %w", err), "")
	}

	return nil
}

func (s *PaymentPlanService) Restore(ctx context.Context, id int) error {
	tx, err := s.repo.BeginTx(ctx)
	if err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to iniciate transaction while restoring payment plan: %w", err), "")
	}

	defer tx.Rollback()

	defer func(){
		p := recover()
		if p != nil{
			s.logger.Error("panic restoring payment plan", "panic", p)
			panic(p)
		}
	}()
	
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	if paymentPlan == nil{
		return apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}

	
	rows, err := s.repo.Restore(ctx, tx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0{
		return apperrors.NewBusinessError(errors.New("failed to restore payment plan: is already restored or doesn't exist"), "El plan de pago ya está activo o no existe." )
	}

	rows, err = s.paymentRepo.BulkUpdateStatus(ctx, tx, paymentPlan.PaymentsInPlan, "in_payment_plan")
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	err = tx.Commit()
	if err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to commit transaction while cancelling payment plan: %w", err), "")
	}

	return nil
}

func (s *PaymentPlanService) HardDelete(ctx context.Context, id int) error {
	err := s.repo.HardDelete(ctx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *PaymentPlanService) RefreshStatus(ctx context.Context, tx *sqlx.Tx, id int) error {
	installments, err := s.installmentRepo.FindAll(ctx, tx, id)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	status := "completed"
	for _, installment := range installments{
		if installment.PaidAt == nil{
			status = "pending"
			break
		}
	}

	err = s.repo.UpdateStatus(ctx, tx, id, status)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}

func (s *PaymentPlanService) Count(ctx context.Context, companyID int) (int, error) {
	count, err := s.repo.Count(ctx, companyID)
	if err!=nil{
		return 0, apperrors.NewDatabaseError(err, "")
	}

	return count, nil
}