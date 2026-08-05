package paymentplan

import (
	"github.com/LucasBastino/app-sindicato/internal/common/page"
)

type pageData struct {
	PaymentPlan     response
	OverduePayments []overduePaymentOption
	CompanyID       int
	PageContext     page.PageContext
	Errors          map[string]string
}

type tablePageData struct {
	PaymentPlans []tableResponse
	TotalResults int
	EmptyState   page.EmptyState
	CompanyID    int
	PageContext  page.PageContext
}

type overviewPageData struct {
	Groups      []PlanCompanyGroup
	EmptyState  page.EmptyState
	PageContext page.PageContext
}
