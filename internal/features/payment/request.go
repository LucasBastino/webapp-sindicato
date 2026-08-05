package payment

import (
	"strings"

	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type request struct {
	Amount       string `form:"amount"`
	PaidAt       string `form:"paid-at"`
	IsPaid       string `form:"is_paid"`
	Observations string `form:"observations"`
}

func (req *request) trim() {
	req.Amount = strings.TrimSpace(req.Amount)
	req.PaidAt = strings.TrimSpace(req.PaidAt)
	req.IsPaid = strings.TrimSpace(req.IsPaid)
	req.Observations = strings.TrimSpace(req.Observations)
}

func (req request) markedPaid() bool {
	return req.IsPaid == "true"
}

// no chequeo companyID porque el payment no puede cambiar de company
func (req request) validate() map[string]string {
	errorMap := map[string]string{}

	if err := v.ValidateAmount(req.Amount); err != "" {
		errorMap["amount"] = err
	}
	if req.Amount == "" {
		errorMap["amount"] = "Campo requerido."
	}
	if err := v.ValidateObservations(req.Observations); err != "" {
		errorMap["observations"] = err
	}

	if req.markedPaid() {
		if req.PaidAt == "" {
			errorMap["paidAt"] = "Ingresá la fecha de pago."
		} else if err := v.ValidatePaidAt(req.PaidAt); err != "" {
			errorMap["paidAt"] = err
		}
	}

	return errorMap
}
