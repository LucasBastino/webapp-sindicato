package payment

import (
	"time"
)

type Payment struct {
	ID		     		int       	`db:"id_payment"`
	CompanyID 			int       	`db:"id_company"`

	Month        		int   	  	`db:"month"`
	Year         		int       	`db:"year"`
	DueDate				time.Time	`db:"due_date"`
	Amount       		*float32    `db:"amount"`
	IsInPaymentPlan		bool		`db:"is_in_payment_plan"`
	PaidAt		  		*time.Time 	`db:"paid_at"`
	
	Observations 		string    	`db:"observations"`

	UpdatedAt    		time.Time 	`db:"updated_at"`
}

func (p Payment) GetStatus() string{
	if p.IsInPaymentPlan {
		return "En plan de pago"
	}
	
	if p.PaidAt != nil {
		return "Completado"
	}
	if time.Now().After(p.DueDate) {
		return "Vencido"
	}
	return "Pendiente"
}