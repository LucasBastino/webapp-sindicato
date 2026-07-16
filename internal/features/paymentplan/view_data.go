package paymentplan

import (
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	"github.com/LucasBastino/app-sindicato/internal/features/installment"
)

type pageData struct {
	PaymentPlan  response
	Installments []installment.Installment
	PageContext  page.PageContext
	Errors       map[string]string
}

type tablePageData struct {
	PaymentPlans	[]tableResponse
	TotalResults 	int

	PageContext  	page.PageContext
}
