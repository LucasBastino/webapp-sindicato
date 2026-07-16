package payment

type response struct {
	Month  int
	Year   int
	Amount float32
	Status string
	PaidAt string

	Observations string

	UpdatedAt string
}

type gridResponse struct {
	Month  int
	Year   int
	Amount float32
	Status string
	PaidAt string
}
