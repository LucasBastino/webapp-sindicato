package company

import "github.com/LucasBastino/app-sindicato/internal/common/page"

type pageData struct {
	Company          response
	Companies        []response
	WithPaymentTable bool
	PageContext      page.PageContext
	Errors           map[string]string
}

type tablePageData struct {
	Companies    []tableResponse
	TotalResults int
	EmptyState   page.EmptyState

	PageContext  page.PageContext
}

// type companyMembersPageData struct {
// 	Members         []member.tableResponse
// 	CompanyID   	int
// 	TotalResults 	int
// 	FromCompany  bool
// 	PageContext page.PageContext
// }

type optionsPageData struct {
	Companies []optionsResponse
	
	PageContext page.PageContext
}
