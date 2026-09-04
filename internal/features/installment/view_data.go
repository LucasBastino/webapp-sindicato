package installment

import "github.com/LucasBastino/webapp-sindicato/internal/common/page"

type pageData struct {
	Installment        response
	PageContext        page.PageContext
	Errors             map[string]string
	CanEditInstallment bool
	PlanCancelled      bool
}
