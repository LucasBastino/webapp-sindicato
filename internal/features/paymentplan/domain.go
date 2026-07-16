package paymentplan

import "time"

type PaymentPlan struct {
	ID                   	int
	CompanyID         		int

	Amount               	float32
	PaymentsInPlan			string
	NumberOfInstallments 	int
	Status               	string
	FirstDueDate           	time.Time
	LastDueDate				time.Time

	Observations			string
	
	CreatedAt				time.Time
	UpdatedAt				time.Time
}

