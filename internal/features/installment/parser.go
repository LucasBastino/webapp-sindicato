package installment

import (
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

func toModel(req request) (Installment, error) {
	installment := Installment{
		Observations: req.Observations,
	}
	if req.PaidAt == "" {
		return installment, nil
	}
	paidAt, err := v.ParseDMY(req.PaidAt)
	if err != nil {
		return Installment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse paid-at: %w", err), "")
	}
	installment.PaidAt = &paidAt
	return installment, nil
}

// cuando da error el formulario queriendo editar
func mergetoResponse(i Installment, req request) (response, error) {
	status := i.GetStatus()
	if req.PaidAt == "" {
		updatedInstallment := i
		updatedInstallment.PaidAt = nil
		status = updatedInstallment.GetStatus()
	}

	return response{
		ID:                i.ID,
		PaymentPlanID:     i.PaymentPlanID,
		InstallmentNumber: i.InstallmentNumber,
		Amount:            i.Amount,
		Status:            status,
		DueDate:           i.DueDate.Format("02/01/2006"),
		PaidAt:            req.PaidAt,
		IsPaid:            req.PaidAt != "",
		Observations:      req.Observations,
		UpdatedAt:         i.UpdatedAt.Format("02/01/2006"),
	}, nil
}

func toResponse(i Installment) response {
	paidAtStr := ""
	if i.PaidAt != nil {
		paidAtStr = i.PaidAt.Format("02/01/2006")
	}
	return response{
		ID:                i.ID,
		PaymentPlanID:     i.PaymentPlanID,
		InstallmentNumber: i.InstallmentNumber,
		Amount:            i.Amount,
		Status:            i.GetStatus(),
		DueDate:           i.DueDate.Format("02/01/2006"),
		PaidAt:            paidAtStr,
		IsPaid:            i.PaidAt != nil,
		Observations:      i.Observations,
		UpdatedAt:         i.UpdatedAt.Format("02/01/2006"),
	}
}

func toGridResponse(i Installment) gridResponse {
	paidAtStr := ""
	if i.PaidAt != nil {
		paidAtStr = i.PaidAt.Format("02/01/2006")
	}
	return gridResponse{
		ID:                i.ID,
		InstallmentNumber: i.InstallmentNumber,
		Amount:            i.Amount,
		Status:            i.GetStatus(),
		DueDate:           i.DueDate.Format("02/01/2006"),
		PaidAt:            paidAtStr,
	}
}

func toGridResponses(installments []Installment) []gridResponse {
	responses := make([]gridResponse, len(installments))
	for i, inst := range installments {
		responses[i] = toGridResponse(inst)
	}
	return responses
}
