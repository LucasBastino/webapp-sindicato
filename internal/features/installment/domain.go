package installment

import "time"

type Installment struct {
	ID                	int
	PaymentPlanID     	int
	
	InstallmentNumber 	int
	Amount              float32
	DueDate      		time.Time
	PaidAt				*time.Time

	Observations      	string

	UpdatedAt			time.Time
}

type paymentPlanData struct {
	id int
	numberOfInstallments int
	amount float32
	firstDueDate time.Time
}

func (i Installment) GetStatus() string{
	if i.PaidAt != nil {
		return "Completado"
	} else {
		if time.Now().After(i.DueDate) {
			return "Vencido"
		} else{
			return "Pendiente"
		}
	}
}