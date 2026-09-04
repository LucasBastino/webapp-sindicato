package paymentplan

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	pu "github.com/LucasBastino/webapp-sindicato/internal/common/utils/parser"
	"github.com/LucasBastino/webapp-sindicato/internal/features/payment"
	v "github.com/LucasBastino/webapp-sindicato/internal/validation"
)

func toCreateInput(req request) (CreateInput, error) {
	companyID, err := strconv.Atoi(req.CompanyID)
	if err != nil {
		return CreateInput{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse company-id: %w", err), "")
	}
	numberOfInstallments, err := strconv.Atoi(req.NumberOfInstallments)
	if err != nil {
		return CreateInput{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse number-of-installments: %w", err), "")
	}
	firstDueDate, err := v.ParseDMY(req.FirstDueDate)
	if err != nil {
		return CreateInput{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse first-due-date: %w", err), "")
	}
	paymentIDs, err := req.parsePaymentIDs()
	if err != nil {
		return CreateInput{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment-ids: %w", err), "")
	}
	normalizedAmount := v.NormalizeAmountInput(req.Amount)
	amount64, err := strconv.ParseFloat(normalizedAmount, 32)
	if err != nil {
		return CreateInput{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse amount: %w", err), "")
	}

	return CreateInput{
		CompanyID:            companyID,
		PaymentIDs:           paymentIDs,
		Amount:               float32(amount64),
		NumberOfInstallments: numberOfInstallments,
		FirstDueDate:         firstDueDate,
		Observations:         req.Observations,
	}, nil
}

func toResponseFromRequest(req request) response {
	amount, _ := strconv.ParseFloat(v.NormalizeAmountInput(req.Amount), 32)
	numberOfInstallments, _ := strconv.Atoi(req.NumberOfInstallments)
	selected := map[int]bool{}
	for _, id := range req.PaymentIDs {
		n, err := strconv.Atoi(id)
		if err == nil {
			selected[n] = true
		}
	}
	return response{
		Amount:               float32(amount),
		NumberOfInstallments: numberOfInstallments,
		FirstDueDate:         req.FirstDueDate,
		Observations:         req.Observations,
		SelectedPaymentIDs:   selected,
	}
}

func mergetoResponse(detail PaymentPlanDetail, req request) (response, error) {
	res := toResponse(detail)
	res.Observations = pu.MergeField(detail.Observations, req.Observations)
	return res, nil
}

func toResponse(detail PaymentPlanDetail) response {
	p := detail.PaymentPlan
	res := response{
		ID:                   p.ID,
		CompanyID:            p.CompanyID,
		CompanyName:          p.CompanyName,
		PaymentsInPlan:       p.PaymentsInPlan,
		OriginalDebt:         detail.OriginalDebt,
		Amount:               p.Amount,
		Status:               p.Status,
		StatusLabel:          statusLabel(p.Status),
		NumberOfInstallments: p.NumberOfInstallments,
		FirstDueDate:         p.FirstDueDate.Format("02/01/2006"),
		LastDueDate:          p.LastDueDate.Format("02/01/2006"),
		Observations:         p.Observations,
		CreatedAt:            p.CreatedAt.Format("02/01/2006"),
		UpdatedAt:            p.UpdatedAt.Format("02/01/2006"),
		IncludedPayments:     toIncludedPayments(detail.IncludedPayments),
		Installments:         make([]installmentResponse, 0, len(detail.Installments)),
	}

	var nextDue *time.Time
	for _, inst := range detail.Installments {
		item := installmentResponse{
			ID:                inst.ID,
			InstallmentNumber: inst.InstallmentNumber,
			Amount:            inst.Amount,
			DueDate:           inst.DueDate.Format("02/01/2006"),
		}
		switch inst.Status {
		case "Completado":
			item.Status = "Pagada"
			item.StatusClass = "green"
			res.PaidCount++
			res.TotalPaid += inst.Amount
			if inst.PaidAt != nil {
				item.PaidAt = inst.PaidAt.Format("02/01/2006")
				item.DateLabel = "Se pagó el " + item.PaidAt
			}
		case "Vencido":
			item.Status = "Vencida"
			item.StatusClass = "red"
			item.DateLabel = "Venció el " + item.DueDate
			res.OverdueCount++
		default:
			item.Status = "Pendiente"
			item.StatusClass = "gray"
			item.DateLabel = "Vence el " + item.DueDate
			res.PendingCount++
			if nextDue == nil || inst.DueDate.Before(*nextDue) {
				due := inst.DueDate
				nextDue = &due
			}
		}
		res.Installments = append(res.Installments, item)
	}
	if nextDue != nil {
		res.NextDueDate = nextDue.Format("02/01/2006")
		res.HasNextDueDate = true
	}
	return res
}

func toIncludedPayments(payments []payment.Payment) []includedPaymentResponse {
	sorted := make([]payment.Payment, len(payments))
	copy(sorted, payments)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Year != sorted[j].Year {
			return sorted[i].Year < sorted[j].Year
		}
		return sorted[i].Month < sorted[j].Month
	})

	out := make([]includedPaymentResponse, 0, len(sorted))
	for _, p := range sorted {
		amount := float32(0)
		if p.Amount != nil {
			amount = *p.Amount
		}
		out = append(out, includedPaymentResponse{
			MonthName: monthName(p.Month),
			Year:      p.Year,
			Amount:    amount,
		})
	}
	return out
}

func toTableResponse(p PaymentPlan) tableResponse {
	return tableResponse{
		ID:                   p.ID,
		Amount:               p.Amount,
		Status:               p.Status,
		StatusLabel:          statusLabel(p.Status),
		StatusBadgeClass:     statusBadgeClass(p.Status),
		NumberOfInstallments: p.NumberOfInstallments,
		FirstDueDate:         p.FirstDueDate.Format("02/01/2006"),
		LastDueDate:          p.LastDueDate.Format("02/01/2006"),
	}
}

func statusBadgeClass(status string) string {
	switch status {
	case "completed":
		return "completed"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}

func toTableResponses(paymentPlans []PaymentPlan) []tableResponse {
	responses := make([]tableResponse, len(paymentPlans))
	for i, plan := range paymentPlans {
		responses[i] = toTableResponse(plan)
	}
	return responses
}

func toOverdueOptions(payments []payment.Payment) []overduePaymentOption {
	opts := make([]overduePaymentOption, 0, len(payments))
	for _, p := range payments {
		amount := float32(0)
		if p.Amount != nil {
			amount = *p.Amount
		}
		opts = append(opts, overduePaymentOption{
			ID:        p.ID,
			MonthName: monthName(p.Month),
			Year:      p.Year,
			Amount:    amount,
			DueDate:   p.DueDate.Format("02/01/2006"),
		})
	}
	return opts
}

var monthNames = []string{
	"", "Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio",
	"Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre",
}

func monthName(month int) string {
	if month < 1 || month > 12 {
		return ""
	}
	return monthNames[month]
}
