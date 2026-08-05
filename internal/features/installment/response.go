package installment

type response struct {
	ID                int
	PaymentPlanID     int
	InstallmentNumber int
	Amount            float32
	Status            string
	DueDate           string
	PaidAt            string
	IsPaid            bool
	Observations      string
	UpdatedAt         string
}

type gridResponse struct {
	ID                int
	InstallmentNumber int
	Amount            float32
	Status            string
	DueDate           string
	PaidAt            string
}
