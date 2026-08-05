package payment

import "github.com/LucasBastino/app-sindicato/internal/common/page"

type pageData struct {
	Payment     response
	CompanyID   int
	CompanyName string
	PageContext page.PageContext
	Errors      map[string]string
}

type gridPageData struct {
	Payments    []gridResponse
	Stats       gridStats
	CompanyID   int
	CompanyName string
	Year        int
	Years       []int
	PageContext page.PageContext
}

type overduePageData struct {
	Groups      []OverdueCompanyGroup
	EmptyState  page.EmptyState
	PageContext page.PageContext
}
