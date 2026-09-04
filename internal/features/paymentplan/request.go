package paymentplan

import (
	"strconv"
	"strings"

	v "github.com/LucasBastino/webapp-sindicato/internal/validation"
)

type request struct {
	CompanyID            string   `form:"company-id"`
	PaymentIDs           []string `form:"payment-ids"`
	PaymentsInPlan       string   `form:"payments-in-plan"`
	Amount               string   `form:"amount"`
	NumberOfInstallments string   `form:"number-of-installments"`
	FirstDueDate         string   `form:"first-due-date"`
	Observations         string   `form:"observations"`
}

func (req *request) trim() {
	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.PaymentsInPlan = strings.TrimSpace(req.PaymentsInPlan)
	req.Amount = strings.TrimSpace(req.Amount)
	req.NumberOfInstallments = strings.TrimSpace(req.NumberOfInstallments)
	req.FirstDueDate = strings.TrimSpace(req.FirstDueDate)
	req.Observations = strings.TrimSpace(req.Observations)
	for i := range req.PaymentIDs {
		req.PaymentIDs[i] = strings.TrimSpace(req.PaymentIDs[i])
	}
}

func (req request) validate() map[string]string {
	errorMap := map[string]string{}

	if err := v.ValidateCompanyID(req.CompanyID); err != "" {
		errorMap["companyID"] = err
	}
	if len(req.PaymentIDs) == 0 {
		errorMap["payment-ids"] = "Debés seleccionar al menos un aporte vencido."
	}
	if req.Amount == "" {
		errorMap["amount"] = "Campo requerido."
	} else if err := v.ValidateAmount(req.Amount); err != "" {
		errorMap["amount"] = err
	} else {
		amount, err := strconv.ParseFloat(v.NormalizeAmountInput(req.Amount), 32)
		if err != nil || amount <= 0 {
			errorMap["amount"] = "El monto debe ser mayor a cero."
		}
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

func validateObservationsOnly(observations string) string {
	return v.ValidateObservations(observations)
}

func (req request) parsePaymentIDs() ([]int, error) {
	ids := make([]int, 0, len(req.PaymentIDs))
	for _, raw := range req.PaymentIDs {
		if raw == "" {
			continue
		}
		id, err := strconv.Atoi(raw)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
