package installment

import "time"

type Installment struct {
	ID                int        `db:"id_installment"`
	PaymentPlanID     int        `db:"id_payment_plan"`
	InstallmentNumber int        `db:"installment_number"`
	Amount            float32    `db:"amount"`
	DueDate           time.Time  `db:"due_date"`
	PaidAt            *time.Time `db:"paid_at"`
	Observations      string     `db:"observations"`
	UpdatedAt         time.Time  `db:"updated_at"`
}


func (i Installment) GetStatus() string {
	if i.PaidAt != nil {
		return "Completado"
	}
	if time.Now().After(i.DueDate) {
		return "Vencido"
	}
	return "Pendiente"
}
