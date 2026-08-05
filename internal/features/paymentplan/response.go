package paymentplan

type response struct {
	ID                   int
	CompanyID            int
	CompanyName          string
	Installments         []installmentResponse
	IncludedPayments     []includedPaymentResponse
	PaymentsInPlan       string
	OriginalDebt         float32
	Amount               float32
	NumberOfInstallments int
	Status               string
	StatusLabel          string
	FirstDueDate         string
	LastDueDate          string
	Observations         string
	CreatedAt            string
	UpdatedAt            string
	SelectedPaymentIDs   map[int]bool
	PaidCount            int
	PendingCount         int
	OverdueCount         int
	TotalPaid            float32
	NextDueDate          string
	HasNextDueDate       bool
}

type includedPaymentResponse struct {
	MonthName string
	Year      int
	Amount    float32
}

type installmentResponse struct {
	ID                int
	InstallmentNumber int
	Amount            float32
	Status            string
	StatusClass       string
	DueDate           string
	PaidAt            string
	DateLabel         string
}

type tableResponse struct {
	ID                   int
	Amount               float32
	NumberOfInstallments int
	Status               string
	StatusLabel          string
	StatusBadgeClass     string
	FirstDueDate         string
	LastDueDate          string
}

type overduePaymentOption struct {
	ID        int
	MonthName string
	Year      int
	Amount    float32
	DueDate   string
}

type overviewPlanResponse struct {
	ID                   int
	Amount               float32
	NumberOfInstallments int
	Status               string
	StatusKey            string
	FirstDueDate         string
	LastDueDate          string
}

type PlanCompanyGroup struct {
	CompanyID   int
	CompanyName string
	TotalAmount float32
	Plans       []overviewPlanResponse
}
