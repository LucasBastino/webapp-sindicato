package payment

type response struct {
	ID     int
	Month  int
	MonthName string
	Year   int
	Amount float32
	Status string
	PaidAt string
	IsPaid bool

	Observations string
	IsInPaymentPlan bool

	DueDateRFC string
	UpdatedAt string
}

type gridResponse struct {
	ID        int
	Month     int
	MonthName string
	Year      int
	Amount    float32
	Status    string
	PaidAt    string
	DueDate   string
	DateLabel string
}

type gridStats struct {
	CompletedCount  int
	CompletedAmount float32
	InPlanCount     int
	InPlanAmount    float32
	OverdueCount    int
	OverdueAmount   float32
	PendingCount    int
	PendingAmount   float32
}

type overduePaymentResponse struct {
	ID        int
	MonthName string
	Year      int
	Amount    float32
	DueDate   string
}

type OverdueCompanyGroup struct {
	CompanyID   int
	CompanyName string
	TotalAmount float32
	Payments    []overduePaymentResponse
}
