package installment

import (
	"context"
	"errors"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/jmoiron/sqlx"
)

type InstallmentService struct {
	repo                       *InstallmentRepository
	paymentPlanStatusRefresher paymentPlanStatusRefresher
}

func NewInstallmentService(repo *InstallmentRepository) *InstallmentService {
	return &InstallmentService{repo: repo}
}

func (s *InstallmentService) SetPaymentPlanStatusRefresher(refresher paymentPlanStatusRefresher) {
	s.paymentPlanStatusRefresher = refresher
}

func (s *InstallmentService) Get(ctx context.Context, id int) (*Installment, error) {
	installment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	if installment == nil {
		return nil, apperrors.NewNotFoundError(errors.New("installment not found"), "")
	}
	return installment, nil
}

func (s *InstallmentService) GetPaymentPlanStatus(ctx context.Context, paymentPlanID int) (string, error) {
	if s.paymentPlanStatusRefresher == nil {
		return "", apperrors.NewInternalError(errors.New("payment plan status refresher not configured"), "")
	}
	return s.paymentPlanStatusRefresher.GetStatus(ctx, paymentPlanID)
}

func (s *InstallmentService) List(ctx context.Context, paymentPlanID int) ([]Installment, error) {
	installments, err := s.repo.FindAll(ctx, nil, paymentPlanID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	return installments, nil
}

func (s *InstallmentService) Update(ctx context.Context, id int, installment Installment) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if existing == nil {
		return apperrors.NewNotFoundError(errors.New("installment not found"), "")
	}

	if s.paymentPlanStatusRefresher == nil {
		return apperrors.NewInternalError(errors.New("payment plan status refresher not configured"), "")
	}
	planStatus, err := s.paymentPlanStatusRefresher.GetStatus(ctx, existing.PaymentPlanID)
	if err != nil {
		return err
	}
	if planStatus == "cancelled" {
		return apperrors.NewBusinessError(
			errors.New("cannot edit installment of cancelled payment plan"),
			"No se puede editar cuotas de un plan cancelado. Restaurá el plan primero.",
		)
	}

	installment.PaymentPlanID = existing.PaymentPlanID
	err = s.repo.Update(ctx, nil, id, installment)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	return s.paymentPlanStatusRefresher.RefreshStatus(ctx, existing.PaymentPlanID)
}

// UpdateInTx is used when the caller already owns a transaction.
func (s *InstallmentService) UpdateInTx(ctx context.Context, tx *sqlx.Tx, id int, installment Installment) error {
	return s.repo.Update(ctx, tx, id, installment)
}
