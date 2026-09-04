package parent

import "github.com/LucasBastino/webapp-sindicato/internal/common/page"

type memberCardData struct {
	ID           int
	Name         string
	LastName     string
	MemberNumber string
	Dni          string
	Phone        string
	CompanyID    int
	CompanyName  string
}

type tablePageData struct {
	Parents      []response
	TotalResults int
	EmptyState   page.EmptyState
	MemberID     int
	MemberName   string
	Member       memberCardData
	PageContext  page.PageContext
}

type pageData struct {
	Parent      response
	MemberID    int
	PageContext page.PageContext
	Errors      map[string]string
}
