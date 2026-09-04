package company

import "github.com/LucasBastino/webapp-sindicato/internal/common/page"

type pageData struct {
	Company          response
	Companies        []response
	WithPaymentTable bool
	PageContext      page.PageContext
	CompanyNav       page.CompanyNav
	Errors           map[string]string
}

type tablePageData struct {
	Companies    []tableResponse
	TotalResults int
	EmptyState   page.EmptyState
	PageContext  page.PageContext
}

type optionsPageData struct {
	Companies   []optionsResponse
	PageContext page.PageContext
}
