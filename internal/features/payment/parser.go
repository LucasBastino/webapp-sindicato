package payment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	pu "github.com/LucasBastino/webapp-sindicato/internal/common/utils/parser"
	v "github.com/LucasBastino/webapp-sindicato/internal/validation"
)

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

func formatPaidAt(paidAt *time.Time) string {
	if paidAt == nil {
		return ""
	}
	return paidAt.Format("02/01/2006")
}

func amountValue(amount *float32) float32 {
	if amount == nil {
		return 0
	}
	return *amount
}

func toModel(req request) (Payment, error) {
	normalizedAmount := v.NormalizeAmountInput(req.Amount)
	amount, err := strconv.ParseFloat(normalizedAmount, 32)
	if err != nil {
		return Payment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse amount: %w", err), "")
	}
	amt := float32(amount)

	payment := Payment{
		Amount:       &amt,
		PaidAt:       nil,
		Observations: req.Observations,
	}

	if !req.markedPaid() {
		return payment, nil
	}

	paidAt, err := v.ParseDMY(req.PaidAt)
	if err != nil {
		return Payment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse paid-at: %w", err), "")
	}
	payment.PaidAt = &paidAt
	return payment, nil
}

func mergetoResponse(p Payment, req request) (response, error) {
	normalizedAmount := v.NormalizeAmountInput(req.Amount)
	amount, err := strconv.ParseFloat(normalizedAmount, 32)
	if err != nil {
		amount = float64(amountValue(p.Amount))
	}
	updatedPayment := p
	isPaid := req.markedPaid()
	if !isPaid {
		updatedPayment.PaidAt = nil
	} else if req.PaidAt == "" {
		updatedPayment.PaidAt = nil
	} else {
		paidAt, err := v.ParseDMY(req.PaidAt)
		if err != nil {
			updatedPayment.PaidAt = nil
		} else {
			updatedPayment.PaidAt = &paidAt
		}
	}

	paidAtStr := formatPaidAt(updatedPayment.PaidAt)
	if isPaid && req.PaidAt != "" {
		paidAtStr = req.PaidAt
	}

	return response{
		ID:              p.ID,
		Month:           p.Month,
		MonthName:       monthName(p.Month),
		Year:            p.Year,
		Status:          updatedPayment.GetStatus(),
		Amount:          pu.MergeField(amountValue(p.Amount), float32(amount)),
		PaidAt:          paidAtStr,
		IsPaid:          isPaid,
		Observations:    pu.MergeField(p.Observations, req.Observations),
		IsInPaymentPlan: p.IsInPaymentPlan,
		DueDateRFC:      p.DueDate.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format("02/01/2006"),
	}, nil
}

func toResponse(p Payment) response {
	return response{
		ID:              p.ID,
		Month:           p.Month,
		MonthName:       monthName(p.Month),
		Year:            p.Year,
		Status:          p.GetStatus(),
		Amount:          amountValue(p.Amount),
		PaidAt:          formatPaidAt(p.PaidAt),
		IsPaid:          p.PaidAt != nil,
		Observations:    p.Observations,
		IsInPaymentPlan: p.IsInPaymentPlan,
		DueDateRFC:      p.DueDate.Format(time.RFC3339),
		UpdatedAt:       p.UpdatedAt.Format("02/01/2006"),
	}
}

func toGridResponse(p Payment) gridResponse {
	status := p.GetStatus()
	dueDate := p.DueDate.Format("02/01/2006")
	paidAt := formatPaidAt(p.PaidAt)

	dateLabel := dueDate
	switch status {
	case "Completado":
		if paidAt != "" {
			dateLabel = "Se pagó el " + paidAt
		}
	case "Vencido":
		if dueDate != "" {
			dateLabel = "Venció el " + dueDate
		}
	case "Pendiente":
		if dueDate != "" {
			dateLabel = "Vence el " + dueDate
		}
	case "En plan de pago":
		if dueDate != "" {
			dateLabel = dueDate
		}
	}

	return gridResponse{
		ID:        p.ID,
		Month:     p.Month,
		MonthName: monthName(p.Month),
		Year:      p.Year,
		Status:    status,
		Amount:    amountValue(p.Amount),
		PaidAt:    paidAt,
		DueDate:   dueDate,
		DateLabel: dateLabel,
	}
}

func toGridResponses(payments []Payment) []gridResponse {
	responses := make([]gridResponse, len(payments))
	for i, payment := range payments {
		responses[i] = toGridResponse(payment)
	}
	return responses
}

func buildGridStats(payments []Payment) gridStats {
	var stats gridStats
	for _, p := range payments {
		amount := amountValue(p.Amount)
		switch p.GetStatus() {
		case "Completado":
			stats.CompletedCount++
			stats.CompletedAmount += amount
		case "En plan de pago":
			stats.InPlanCount++
			stats.InPlanAmount += amount
		case "Vencido":
			stats.OverdueCount++
			stats.OverdueAmount += amount
		default:
			stats.PendingCount++
			stats.PendingAmount += amount
		}
	}
	return stats
}
