package paymentplan

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/features/installment"
	"github.com/LucasBastino/app-sindicato/internal/features/payment"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/jmoiron/sqlx"
)

type PaymentPlanService struct {
	repo            *PaymentPlanRepository
	paymentRepo     *payment.PaymentRepository
	installmentRepo *installment.InstallmentRepository
	idempotency     *idempotency.IdempotencyService
	logger          logger.Logger
}

func NewPaymentPlanService(
	repo *PaymentPlanRepository,
	paymentRepo *payment.PaymentRepository,
	installmentRepo *installment.InstallmentRepository,
	idempotencyService *idempotency.IdempotencyService,
	logger logger.Logger,
) *PaymentPlanService {
	return &PaymentPlanService{
		repo:            repo,
		paymentRepo:     paymentRepo,
		installmentRepo: installmentRepo,
		idempotency:     idempotencyService,
		logger:          logger,
	}
}

func (s *PaymentPlanService) Get(ctx context.Context, id int) (*PaymentPlan, []installment.Installment, error) {
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, nil, apperrors.NewDatabaseError(err, "")
	}
	if paymentPlan == nil {
		return nil, nil, apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}

	installments, err := s.installmentRepo.FindAll(ctx, nil, id)
	if err != nil {
		return nil, nil, apperrors.NewDatabaseError(err, "")
	}
	return paymentPlan, installments, nil
}

func (s *PaymentPlanService) GetDetail(ctx context.Context, id int) (*PaymentPlanDetail, error) {
	paymentPlan, installments, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	paymentIDs, err := parsePaymentIDs(paymentPlan.PaymentsInPlan)
	if err != nil {
		return nil, apperrors.NewInternalError(fmt.Errorf("invalid payments_in_plan: %w", err), "")
	}
	payments, err := s.paymentRepo.FindByIDs(ctx, paymentIDs)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}

	detail := &PaymentPlanDetail{
		PaymentPlan:      *paymentPlan,
		IncludedPayments: payments,
		Installments:     make([]InstallmentDetail, 0, len(installments)),
	}
	for _, p := range payments {
		if p.Amount != nil {
			detail.OriginalDebt += *p.Amount
		}
	}
	for _, inst := range installments {
		detail.Installments = append(detail.Installments, InstallmentDetail{
			ID:                inst.ID,
			InstallmentNumber: inst.InstallmentNumber,
			Amount:            inst.Amount,
			Status:            inst.GetStatus(),
			DueDate:           inst.DueDate,
			PaidAt:            inst.PaidAt,
		})
	}
	return detail, nil
}

func (s *PaymentPlanService) List(ctx context.Context, companyID int) ([]PaymentPlan, error) {
	paymentPlans, err := s.repo.FindAll(ctx, companyID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	return paymentPlans, nil
}

func (s *PaymentPlanService) CountAll(ctx context.Context) (int, error) {
	count, err := s.repo.CountAll(ctx)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}
	return count, nil
}

func (s *PaymentPlanService) ListAllGrouped(ctx context.Context) ([]PlanCompanyGroup, error) {
	rows, err := s.repo.FindAllWithCompany(ctx)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	if len(rows) == 0 {
		return nil, nil
	}

	byCompany := make(map[int]*PlanCompanyGroup)
	order := make([]int, 0)
	for _, row := range rows {
		group, ok := byCompany[row.CompanyID]
		if !ok {
			group = &PlanCompanyGroup{
				CompanyID:   row.CompanyID,
				CompanyName: row.CompanyName,
				Plans:       make([]overviewPlanResponse, 0),
			}
			byCompany[row.CompanyID] = group
			order = append(order, row.CompanyID)
		}
		group.TotalAmount += row.Amount
		group.Plans = append(group.Plans, overviewPlanResponse{
			ID:                   row.ID,
			Amount:               row.Amount,
			NumberOfInstallments: row.NumberOfInstallments,
			Status:               statusLabel(row.Status),
			StatusKey:            row.Status,
			FirstDueDate:         row.FirstDueDate.Format("02/01/2006"),
			LastDueDate:          row.LastDueDate.Format("02/01/2006"),
		})
	}

	groups := make([]PlanCompanyGroup, 0, len(order))
	for _, id := range order {
		groups = append(groups, *byCompany[id])
	}
	return groups, nil
}

func statusLabel(status string) string {
	switch status {
	case "completed":
		return "Completado"
	case "cancelled":
		return "Cancelado"
	default:
		return "Pendiente"
	}
}

func (s *PaymentPlanService) ListOverdueForCompany(ctx context.Context, companyID int) ([]payment.Payment, error) {
	payments, err := s.paymentRepo.FindOverdueByCompany(ctx, companyID)
	if err != nil {
		return nil, apperrors.NewDatabaseError(err, "")
	}
	return payments, nil
}

type CreateInput struct {
	CompanyID            int
	PaymentIDs           []int
	Amount               float32
	NumberOfInstallments int
	FirstDueDate         time.Time
	Observations         string
}

func (s *PaymentPlanService) Create(ctx context.Context, input CreateInput, idempotencyKey string) (int, error) {
	if len(input.PaymentIDs) == 0 {
		return 0, apperrors.NewBusinessError(errors.New("no payments selected"), "Debés seleccionar al menos un pago vencido.")
	}
	if input.NumberOfInstallments < 1 {
		return 0, apperrors.NewBusinessError(errors.New("invalid installments"), "La cantidad de cuotas no es válida.")
	}
	if input.Amount <= 0 {
		return 0, apperrors.NewBusinessError(errors.New("invalid amount"), "El monto total del plan debe ser mayor a cero.")
	}

	selected, err := s.paymentRepo.FindByIDs(ctx, input.PaymentIDs)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}
	if len(selected) != len(input.PaymentIDs) {
		return 0, apperrors.NewBusinessError(errors.New("some payments not found"), "Algunos pagos seleccionados no existen.")
	}

	idStrs := make([]string, 0, len(selected))
	for _, p := range selected {
		if p.CompanyID != input.CompanyID {
			return 0, apperrors.NewBusinessError(errors.New("payment company mismatch"), "Los pagos deben pertenecer a la misma empresa.")
		}
		if p.PaidAt != nil || p.IsInPaymentPlan || !time.Now().After(p.DueDate) {
			return 0, apperrors.NewBusinessError(errors.New("payment not eligible"), "Solo se pueden incluir pagos vencidos disponibles.")
		}
		idStrs = append(idStrs, strconv.Itoa(p.ID))
	}

	lastDue := addMonthsClamped(input.FirstDueDate, input.NumberOfInstallments-1)
	plan := PaymentPlan{
		CompanyID:            input.CompanyID,
		PaymentsInPlan:       strings.Join(idStrs, ","),
		Amount:               input.Amount,
		NumberOfInstallments: input.NumberOfInstallments,
		Status:               "pending",
		FirstDueDate:         input.FirstDueDate,
		LastDueDate:          lastDue,
		Observations:         input.Observations,
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx creating payment plan: %w", err), "")
	}
	defer tx.Rollback()

	defer func() {
		if p := recover(); p != nil {
			s.logger.Error("panic creating payment plan", "panic", p)
			panic(p)
		}
	}()

	id, err := s.repo.Insert(ctx, tx, plan)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}

	installments := buildInstallments(id, input.Amount, input.NumberOfInstallments, input.FirstDueDate)
	rows, err := s.installmentRepo.BulkInsert(ctx, tx, installments)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}
	if rows != len(installments) {
		return 0, apperrors.NewBusinessError(errors.New("installment insert mismatch"), "No se pudieron crear todas las cuotas del plan.")
	}

	_, err = s.paymentRepo.BulkUpdateIsInPaymentPlan(ctx, tx, input.PaymentIDs, true)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}

	if err := s.idempotency.UpdateResource(ctx, tx, idempotencyKey, "payment_plan", id); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, apperrors.NewDatabaseError(fmt.Errorf("failed to commit payment plan create: %w", err), "")
	}
	return id, nil
}

func buildInstallments(planID int, total float32, n int, firstDue time.Time) []installment.Installment {
	installments := make([]installment.Installment, 0, n)
	base := float32(int((total/float32(n))*100+0.5)) / 100
	assigned := float32(0)
	for i := 0; i < n; i++ {
		amt := base
		if i == n-1 {
			amt = float32(int((total-assigned)*100+0.5)) / 100
		}
		assigned += amt
		installments = append(installments, installment.Installment{
			PaymentPlanID:     planID,
			InstallmentNumber: i + 1,
			Amount:            amt,
			DueDate:           addMonthsClamped(firstDue, i),
		})
	}
	return installments
}

func addMonthsClamped(t time.Time, months int) time.Time {
	year, month, day := t.Date()
	firstOf := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, months, 0)
	lastDay := time.Date(firstOf.Year(), firstOf.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
	if day > lastDay {
		day = lastDay
	}
	return time.Date(firstOf.Year(), firstOf.Month(), day, 0, 0, 0, 0, time.UTC)
}

func parsePaymentIDs(csv string) ([]int, error) {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil, nil
	}
	parts := strings.Split(csv, ",")
	ids := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *PaymentPlanService) Update(ctx context.Context, id int, paymentPlan PaymentPlan) error {
	err := s.repo.Update(ctx, id, paymentPlan)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	return nil
}

func (s *PaymentPlanService) Cancel(ctx context.Context, id int) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx cancelling payment plan: %w", err), "")
	}
	defer tx.Rollback()

	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if paymentPlan == nil {
		return apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}

	rows, err := s.repo.Cancel(ctx, tx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0 {
		return apperrors.NewBusinessError(errors.New("failed to cancel payment plan"), "El plan de pago ya está cancelado o no se puede cancelar.")
	}

	ids, err := parsePaymentIDs(paymentPlan.PaymentsInPlan)
	if err != nil {
		return apperrors.NewInternalError(fmt.Errorf("invalid payments_in_plan: %w", err), "")
	}
	_, err = s.paymentRepo.BulkUpdateIsInPaymentPlan(ctx, tx, ids, false)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to commit cancel payment plan: %w", err), "")
	}
	return nil
}

func (s *PaymentPlanService) Restore(ctx context.Context, id int) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx restoring payment plan: %w", err), "")
	}
	defer tx.Rollback()

	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if paymentPlan == nil {
		return apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}

	ids, err := parsePaymentIDs(paymentPlan.PaymentsInPlan)
	if err != nil {
		return apperrors.NewInternalError(fmt.Errorf("invalid payments_in_plan: %w", err), "")
	}
	payments, err := s.paymentRepo.FindByIDs(ctx, ids)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if len(payments) != len(ids) {
		return apperrors.NewBusinessError(
			errors.New("payments not restorable"),
			"No se puede restaurar el plan: algún pago del plan ya no existe.",
		)
	}
	now := time.Now()
	for _, p := range payments {
		if p.PaidAt != nil || p.IsInPaymentPlan || !now.After(p.DueDate) {
			return apperrors.NewBusinessError(
				errors.New("payments not restorable"),
				"No se puede restaurar el plan: todos los pagos incluidos deben seguir vencidos y fuera de otro plan.",
			)
		}
	}

	rows, err := s.repo.Restore(ctx, tx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if rows == 0 {
		return apperrors.NewBusinessError(errors.New("failed to restore payment plan"), "El plan de pago ya está activo o no se puede restaurar.")
	}

	_, err = s.paymentRepo.BulkUpdateIsInPaymentPlan(ctx, tx, ids, true)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to commit restore payment plan: %w", err), "")
	}
	return nil
}

func (s *PaymentPlanService) HardDelete(ctx context.Context, id int) error {
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if paymentPlan == nil {
		return apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}
	if paymentPlan.Status != "cancelled" {
		return apperrors.NewBusinessError(
			errors.New("payment plan not cancelled"),
			"Solo se pueden eliminar planes de pago cancelados.",
		)
	}

	err = s.repo.HardDelete(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	return nil
}

// RefreshStatus implements installment.paymentPlanStatusRefresher.
func (s *PaymentPlanService) RefreshStatus(ctx context.Context, id int) error {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to begin tx refreshing payment plan: %w", err), "")
	}
	defer tx.Rollback()

	if err := s.refreshStatusTx(ctx, tx, id); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to commit refresh payment plan: %w", err), "")
	}
	return nil
}

// GetStatus implements installment.paymentPlanStatusRefresher.
func (s *PaymentPlanService) GetStatus(ctx context.Context, id int) (string, error) {
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return "", apperrors.NewDatabaseError(err, "")
	}
	if paymentPlan == nil {
		return "", apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}
	return paymentPlan.Status, nil
}

func (s *PaymentPlanService) refreshStatusTx(ctx context.Context, tx *sqlx.Tx, id int) error {
	paymentPlan, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	if paymentPlan == nil {
		return apperrors.NewNotFoundError(errors.New("payment plan not found"), "")
	}
	if paymentPlan.Status == "cancelled" {
		return nil
	}

	installments, err := s.installmentRepo.FindAll(ctx, tx, id)
	if err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	status := "completed"
	var lastPaid *time.Time
	for _, inst := range installments {
		if inst.PaidAt == nil {
			status = "pending"
			break
		}
		if lastPaid == nil || inst.PaidAt.After(*lastPaid) {
			lastPaid = inst.PaidAt
		}
	}

	if err := s.repo.UpdateStatus(ctx, tx, id, status); err != nil {
		return apperrors.NewDatabaseError(err, "")
	}

	if status == "completed" {
		ids, err := parsePaymentIDs(paymentPlan.PaymentsInPlan)
		if err != nil {
			return apperrors.NewInternalError(fmt.Errorf("invalid payments_in_plan: %w", err), "")
		}
		paidAt := time.Now()
		if lastPaid != nil {
			paidAt = *lastPaid
		}
		_, err = s.paymentRepo.BulkMarkCompleted(ctx, tx, ids, paidAt)
		if err != nil {
			return apperrors.NewDatabaseError(err, "")
		}
	}
	return nil
}

func (s *PaymentPlanService) Count(ctx context.Context, companyID int) (int, error) {
	count, err := s.repo.Count(ctx, companyID)
	if err != nil {
		return 0, apperrors.NewDatabaseError(err, "")
	}
	return count, nil
}
