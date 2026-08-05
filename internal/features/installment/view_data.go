package installment

import "github.com/LucasBastino/app-sindicato/internal/common/page"

type pageData struct {
	Installment        response
	PageContext        page.PageContext
	Errors             map[string]string
	CanEditInstallment bool
	PlanCancelled      bool
}
