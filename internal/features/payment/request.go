package payment

import (
	"strings"

	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type request struct {
	// Month     	    string `form:"month"`
	// Year      	    string `form:"year"`
	Amount     	  	string `form:"amount"`
	PaidAt 			string `form:"paid-at"`
	
	Observations 	string `form:"observations"`
}


func (req *request) trim() {
	// req.Month = strings.TrimSpace(req.Month)
	// req.Year = strings.TrimSpace(req.Year)
	// req.Status = strings.TrimSpace(req.Status)
	req.Amount = strings.TrimSpace(req.Amount)
	req.PaidAt = strings.TrimSpace(req.PaidAt)
	req.Observations = strings.TrimSpace(req.Observations)
}

// no chequeo companyID porque el payment no puede cambiar de company
func (req request) validate() map[string]string {

	errorMap := map[string]string{}

	// if err := v.ValidateMonth(req.Month); err != "" {
	// 	errorMap["month"] = err
	// }
	// if err := v.ValidateYear(req.Year); err != "" {
	// 	errorMap["year"] = err
	// }
	// if err := v.ValidateIsPaid(req.IsPaid); err != "" {
	// 	errorMap["is-paid"] = err
	// }
	// if err := v.ValidateStatus(req.Status); err != "" {
	// 	errorMap["status"] = err
	// }
	if err := v.ValidateAmount(req.Amount); err != "" {
		errorMap["amount"] = err
	}
	if err := v.ValidatePaidAt(req.PaidAt); err != "" {
		errorMap["paidAt"] = err
	}
	if err := v.ValidateObservations(req.Observations); err != "" {
		errorMap["observations"] = err
	}
	return errorMap
}


