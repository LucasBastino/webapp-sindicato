package installment

import (
	"fmt"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	pu "github.com/LucasBastino/app-sindicato/internal/common/utils/parser"
)

func toModel(req request) (Installment, error) {
	paidAt, err := time.Parse("02/01/2006", req.PaidAt)
	if err!=nil{
		return Installment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse paid-at: %w", err), "")
	}

    return Installment{
        PaidAt:  		&paidAt,
        Observations: 	req.Observations,
    }, nil
}

// cuando da error el formulario queriendo editar
func mergetoResponse(i Installment, req request) (response, error) {
	updatedInstallment := i
	if req.PaidAt == ""{
		updatedInstallment.PaidAt = nil
	} else {
		paidAt, err := time.Parse("02/01/2006", req.PaidAt)
		if err!=nil{
			return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse paid-at: %w", err), "")
		}
		updatedInstallment.PaidAt = &paidAt
	}
	return response{
		Amount:       	i.Amount,
		Status:			updatedInstallment.GetStatus(),
		PaidAt:  		pu.MergeField(i.PaidAt.Format("02/01/2006"), req.PaidAt),
		Observations: 	pu.MergeField(i.Observations, req.Observations),
		UpdatedAt:    	i.UpdatedAt.Format("02/01/2006"),
	}, nil
}

func toResponse(i Installment) response {
    return response{
		Amount:     	i.Amount,
        Status:       	i.GetStatus(),
        PaidAt:  		i.PaidAt.Format("02/01/2006"),
        Observations: 	i.Observations,
        UpdatedAt:   	i.UpdatedAt.Format("02/01/2006"),
    }
}

// func toResponses(installments []Installment) []response{
// 	Responses := make([]response, len(installments))
// 	for i, installment := range installments{
// 		Responses[i] = toResponse(installment)
// 	}
// 	return Responses
// }

func toGridResponse(i Installment) gridResponse {
	return gridResponse{
		Amount: i.Amount,
		Status: i.GetStatus(),
		PaidAt:	i.PaidAt.Format("02/01/2006"),
	}
}

func toGridResponses(installments []Installment) []gridResponse{
	responses := make([]gridResponse, len(installments))
	for i, installment := range installments{
		responses[i] = toGridResponse(installment)
	}
	return responses
}