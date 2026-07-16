package installment

import (
	"context"
	"errors"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
)

type InstallmentService struct {
	repo *InstallmentRepository
	
	paymentPlanStatusRefresher paymentPlanStatusRefresher
}

func NewInstallmentService(repo *InstallmentRepository) *InstallmentService {
	return &InstallmentService{
		repo: repo,
	}
}

func (s *InstallmentService) Get(ctx context.Context, id int) (*Installment, error) {
	installment, err := s.repo.FindByID(ctx, id)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	if installment == nil{
		return nil, apperrors.NewNotFoundError(errors.New("installment not found"), "")
	}

	return installment, nil
}

func (s *InstallmentService) List(ctx context.Context, paymentPlanID int) ([]Installment, error) {
	installments, err := s.repo.FindAll(ctx, nil, paymentPlanID)
	if err!=nil{
		return nil, apperrors.NewDatabaseError(err, "")
	}

	return installments, nil
}

func (s *InstallmentService) Create(ctx context.Context, paymentPlan paymentPlanData) error {
	installments := buildInstallments(paymentPlan)

	rows, err := s.repo.BulkInsert(ctx, installments)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	if rows != len(installments){
		return apperrors.NewBusinessError(errors.New("failed to bulk insert all installments: rows affected doesn't match with installments count"), "")
	}

	return nil
}

func (s *InstallmentService) Update(ctx context.Context, id int, installment Installment) error {
	err := s.repo.Update(ctx, id, installment)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return s.paymentPlanStatusRefresher.RefreshStatus(ctx, installment.PaymentPlanID)
}

func (s *InstallmentService) CheckStatusDaily(ctx context.Context) error {
	err := s.repo.CheckStatusDaily(ctx)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

	return nil
}


func buildInstallments(paymentPlan paymentPlanData) []Installment {
	var installments []Installment

	for i:=0; i < paymentPlan.numberOfInstallments; i++ {
		baseDay := paymentPlan.firstDueDate.Day()
		month := paymentPlan.firstDueDate.Month()
		year := paymentPlan.firstDueDate.Year()
		lastDay := daysInMonth(year, month)
		day := baseDay
		if baseDay > lastDay{
			day = lastDay
		}
		amount := paymentPlan.amount / float32(paymentPlan.numberOfInstallments)
		dueDate := time.Date(year, month, day, 0, 0, 0, 0 , time.UTC)
		installments = append(installments, Installment{
			PaymentPlanID: paymentPlan.id,
			InstallmentNumber: i+1 ,
			Amount: amount,
			DueDate: dueDate,
		})
	}
	return installments
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}