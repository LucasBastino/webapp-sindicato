package member



import "github.com/LucasBastino/app-sindicato/internal/common/page"



type pageData struct {

	Member      response

	PageContext page.PageContext

	Errors      map[string]string

}



type tablePageData struct {
	Members      []tableResponse
	TotalResults int
	EmptyState   page.EmptyState
	CompanyID    int
	CompanyName  string
	PageContext  page.PageContext
}


