package paymentplan

import (
	"time"

	"github.com/LucasBastino/app-sindicato/internal/features/payment"
)

type PaymentPlan struct {
	ID                   int       `db:"id_payment_plan"`
	CompanyID            int       `db:"id_company"`
	CompanyName          string    `db:"company_name"`
	Amount               float32   `db:"amount"`
	PaymentsInPlan       string    `db:"payments_in_plan"`
	NumberOfInstallments int       `db:"number_of_installments"`
	Status               string    `db:"status"`
	FirstDueDate         time.Time `db:"first_due_date"`
	LastDueDate          time.Time `db:"last_due_date"`
	Observations         string    `db:"observations"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}

type PaymentPlanDetail struct {
	PaymentPlan
	OriginalDebt     float32
	IncludedPayments []payment.Payment
	Installments     []InstallmentDetail
}

type InstallmentDetail struct {
	ID                int
	InstallmentNumber int
	Amount            float32
	Status            string
	DueDate           time.Time
	PaidAt            *time.Time
}
