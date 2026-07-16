package paymentplan

import (
	"strings"

	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

type request struct {
	PaymentsInPlan       string `form:"payments-in-plan"`
	Amount               string `form:"amount"`
	NumberOfInstallments string `form:"number-of-installments"`
	FirstDueDate         string `form:"first-due-date"`
	
	Observations         string `form:"observations"`
}

func (req *request) trim() {
	req.PaymentsInPlan = strings.TrimSpace(req.PaymentsInPlan)
	req.Amount = strings.TrimSpace(req.Amount)
	req.NumberOfInstallments = strings.TrimSpace(req.NumberOfInstallments)
	req.FirstDueDate = strings.TrimSpace(req.FirstDueDate)
	req.Observations = strings.TrimSpace(req.Observations)
}

// no chequeo companyID porque el payment no puede cambiar de company
func (req request) validate() map[string]string {

	errorMap := map[string]string{}

	if err := v.ValidateAmount(req.Amount); err != "" {
		errorMap["amount"] = err
	}
	if err := v.ValidateNumberOfInstallments(req.NumberOfInstallments); err != "" {
		errorMap["number-of-installments"] = err
	}
	if err := v.ValidateFirstDueDate(req.FirstDueDate); err != "" {
		errorMap["first-due-date"] = err
	}
	if err := v.ValidateObservations(req.Observations); err != "" {
		errorMap["observations"] = err
	}
	return errorMap
}
