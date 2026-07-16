package installment

type response struct {
	NumberOfInstallment int

	Amount float32
	Status string
	PaidAt string

	Observations string

	UpdatedAt string
}

type gridResponse struct {
	NumberOfInstallment int

	Amount float32
	Status string
	PaidAt string
}
