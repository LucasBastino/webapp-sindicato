package paymentplan

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	pu "github.com/LucasBastino/app-sindicato/internal/common/utils/parser"
	"github.com/LucasBastino/app-sindicato/internal/features/installment"
)

func toModel(req request) (PaymentPlan, error) {
	// paymentsInPlan, err := strconv.Atoi(req.PaymentsInPlan)
	// if err!=nil{
	// 	return PaymentPlan{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse payments-in-plan: %w", err), "")
	// }
	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err!=nil{
		return PaymentPlan{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse amount: %w", err), "")
	}
	numberOfInstallments, err := strconv.Atoi(req.NumberOfInstallments)
	if err!=nil{
		return PaymentPlan{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse number-of-installments: %w", err), "")
	}
	firstDueDate, err := time.Parse("02/01/2006", req.FirstDueDate)
	if err!=nil{
		return PaymentPlan{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse first-due-date: %w", err), "")
	}

    return PaymentPlan{
        PaymentsInPlan: 		req.PaymentsInPlan,
		Amount: 				float32(amount),
		NumberOfInstallments: 	numberOfInstallments,
		FirstDueDate: 			firstDueDate,
        Observations: 			req.Observations,
    }, nil
}

func toResponseFromRequest(req request) (response, error) {
	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err!=nil{
		return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse amount: %w", err), "")
	}
	numberOfInstallments, err := strconv.Atoi(req.NumberOfInstallments)
	if err!=nil{
		return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse number of installments: %w", err), "")
	}

    return response{
		PaymentsInPlan:			req.PaymentsInPlan,
		Amount:       			float32(amount),
		NumberOfInstallments:	numberOfInstallments,
		FirstDueDate:  			req.FirstDueDate,
		Observations: 			req.Observations,
    }, nil
}

// cuando da error el formulario queriendo editar
func mergetoResponse(p PaymentPlan, req request) (response, error) {
	return response{
		PaymentsInPlan:			p.PaymentsInPlan,
		Amount:       			p.Amount,
		Status:					p.Status,
		NumberOfInstallments:	p.NumberOfInstallments,
		FirstDueDate:  			p.FirstDueDate.Format("02/01/2006"),
		LastDueDate:  			p.LastDueDate.Format("02/01/2006"),
		Observations: 			pu.MergeField(p.Observations, req.Observations),
		CreatedAt:    			p.CreatedAt.Format("02/01/2006"),
		UpdatedAt:    			p.UpdatedAt.Format("02/01/2006"),
	}, nil
}

func toResponse(p PaymentPlan, installments []installment.Installment) response {
    return response{
		Installments:			installments,
		PaymentsInPlan:			p.PaymentsInPlan,
		Amount:       			p.Amount,
		Status:					p.Status,
		NumberOfInstallments:	p.NumberOfInstallments,
		FirstDueDate:  			p.FirstDueDate.Format("02/01/2006"),
		LastDueDate:  			p.LastDueDate.Format("02/01/2006"),
		Observations: 			p.Observations,
		CreatedAt:    			p.CreatedAt.Format("02/01/2006"),
		UpdatedAt:    			p.UpdatedAt.Format("02/01/2006"),
    }
}

// func toResponses(paymentPlans []PaymentPlan) []response{
// 	Responses := make([]response, len(paymentPlans))
// 	for i, PaymentPlan := range paymentPlans{
// 		Responses[i] = toResponse(PaymentPlan)
// 	}
// 	return Responses
// }

func toTableResponse(p PaymentPlan) tableResponse {
	return tableResponse{
		Amount:       			p.Amount,
		Status:					p.Status,
		NumberOfInstallments:	p.NumberOfInstallments,
		FirstDueDate:  			p.FirstDueDate.Format("02/01/2006"),
		LastDueDate:  			p.LastDueDate.Format("02/01/2006"),
	}
}

func toTableResponses(paymentPlans []PaymentPlan) []tableResponse{
	responses := make([]tableResponse, len(paymentPlans))
	for i, PaymentPlan := range paymentPlans{
		responses[i] = toTableResponse(PaymentPlan)
	}
	return responses
}