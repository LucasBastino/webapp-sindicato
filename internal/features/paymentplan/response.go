package paymentplan

import "github.com/LucasBastino/app-sindicato/internal/features/installment"

type response struct {
	Installments         []installment.Installment
	
	PaymentsInPlan       string
	Amount               float32
	NumberOfInstallments int
	Status               string
	FirstDueDate         string
	LastDueDate          string

	Observations         string

	CreatedAt            string
	UpdatedAt            string
}

type tableResponse struct {
	Amount               float32
	NumberOfInstallments int
	Status               string
	FirstDueDate         string
	LastDueDate          string
}
