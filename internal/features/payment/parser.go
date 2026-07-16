package payment

import (
	"fmt"
	"strconv"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	pu "github.com/LucasBastino/app-sindicato/internal/common/utils/parser"
)

func toModel(req request) (Payment, error) {
	// monthInt, err := strconv.Atoi(req.Amount)
	// if err!=nil{
	// 	return Payment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse month: %w", err), "")
	// }
	// yearInt, err := strconv.Atoi(req.Amount)
	// if err!=nil{
	// 	return Payment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse year: %w", err), "")
	// }
	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err!=nil{
		return Payment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse amount: %w", err), "")
	}
	paidAt, err := time.Parse("02/01/2006", req.PaidAt)
	if err!=nil{
		return Payment{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse paid-at: %w", err), "")
	}

    return Payment{
        // Month:        	monthInt,
        // Year:         	yearInt,
        // Status:  		req.Status,
        Amount:       	float32(amount),
        PaidAt:  		&paidAt,
        Observations: 	req.Observations,
    }, nil
}

// cuando da error el formulario queriendo editar
func mergetoResponse(p Payment, req request) (response, error) {

	amount, err := strconv.ParseFloat(req.Amount, 32)
	if err!=nil{
		return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse amount: %w", err), "")
	}
	updatedPayment := p
	if req.PaidAt == ""{
		updatedPayment.PaidAt = nil
	} else {
		paidAt, err := time.Parse("02/01/2006", req.PaidAt)
		if err!=nil{
			return response{}, apperrors.NewBadRequestError(fmt.Errorf("failed to parse paid-at: %w", err), "")
		}
		updatedPayment.PaidAt = &paidAt
	}

	return response{
		Month:       	p.Month,
		Year:         	p.Year,
		Status:			updatedPayment.GetStatus(),
		Amount:       	pu.MergeField(p.Amount, float32(amount)),
		PaidAt:  		pu.MergeField(p.PaidAt.Format("02/01/2006"), req.PaidAt),
		Observations: 	pu.MergeField(p.Observations, req.Observations),
		UpdatedAt:    	p.UpdatedAt.Format("02/01/2006"),
	}, nil
}

func toResponse(p Payment) response {
    return response{
        Month:        p.Month,
        Year:         p.Year,
        Status:       p.GetStatus(),
        Amount:       p.Amount,
        PaidAt:  	  p.PaidAt.Format("02/01/2006"),
        Observations: p.Observations,
        UpdatedAt:    p.UpdatedAt.Format("02/01/2006"),
    }
}

// func toResponses(payments []Payment) []response{
// 	Responses := make([]response, len(payments))
// 	for i, paymentModel := range payments{
// 		Responses[i] = toResponse(paymentModel)
// 	}
// 	return Responses
// }

func toGridResponse(p Payment) gridResponse {
	return gridResponse{
		Month:       p.Month,
		Year:        p.Year,
		Status:      p.GetStatus(),
		Amount:      p.Amount,
		PaidAt: 	 p.PaidAt.Format("02/01/2006"),
	}
}

func toGridResponses(payments []Payment) []gridResponse{
	responses := make([]gridResponse, len(payments))
	for i, payment := range payments{
		responses[i] = toGridResponse(payment)
	}
	return responses
}

/* 
type PaymentParser struct{}

func (parser PaymentParser) ParseModel(c *fiber.Ctx) (Payment, error) {
	p := Payment{}
	p.Month = strings.TrimSpace(c.FormValue("month"))
	p.Year = strings.TrimSpace(c.FormValue("year"))
	status, err := strconv.ParseBool(strings.TrimSpace(c.FormValue("status")))
	if err != nil {
		customError.InternalError.Msg = err.Error()
		return Payment{}, errorHandler.HandleError(c, customError.InternalError)
	}
	p.Status = status
	AmountStr := strings.TrimSpace(c.FormValue("amount"))
	if AmountStr != "" {
		Amount, err := strconv.Atoi(AmountStr)
		if err != nil {
			customError.StrConvError.Msg = err.Error()
			return Payment{}, customError.StrConvError
		}
		p.Amount = Amount
	}
	p.PaidAt = strings.TrimSpace(c.FormValue("paid-at"))
	p.Observations = strings.TrimSpace(c.FormValue("observations"))
	CompanyIDStr := strings.TrimSpace(c.FormValue("id-company"))
	CompanyID, err := strconv.Atoi(CompanyIDStr)
	if err != nil {
		customError.StrConvError.Msg = err.Error()
		return Payment{}, err
	}
	p.CompanyID = CompanyID

	return p, nil
}
 */



// // cuando da error el formulario en el frontend
// func toResponseFromRequest(req request) response {
// 	// var status bool
// 	// if req.Status == "true"{
// 	// 	status = true
// 	// } else {
// 	// 	status = false
// 	// }
// 	// ya estan validados, no hace falta chequear el error
// 	amountInt, _ := strconv.Atoi(req.Amount)
// 	companyIDInt, _ := strconv.Atoi(req.CompanyID)
// 	monthInt, _ := strconv.Atoi(req.Month)
// 	yearInt, _ := strconv.Atoi(req.Year)

//     return response{
//         Month:       	monthInt,
//         Year:         	yearInt,
//         Status:  status,
//         Amount:      	amountInt,
//         PaidAt:  	req.PaidAt,
//         Observations: 	req.Observations,
//         CompanyID: 	companyIDInt,
// 	}
// }