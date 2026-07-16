package installment

import (
	"strings"

	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type request struct {
	PaidAt 			string `form:"paid-at"`
	Observations 	string `form:"observations"`
}


func (req *request) trim() {
	req.PaidAt = strings.TrimSpace(req.PaidAt)
	req.Observations = strings.TrimSpace(req.Observations)
}

// no chequeo paymentPlanID porque el payment no puede cambiar de company
func (req request) validate() map[string]string {

	errorMap := map[string]string{}

	if err := v.ValidatePaidAt(req.PaidAt); err != "" {
		errorMap["paidAt"] = err
	}
	if err := v.ValidateObservations(req.Observations); err != "" {
		errorMap["observations"] = err
	}
	return errorMap
}


