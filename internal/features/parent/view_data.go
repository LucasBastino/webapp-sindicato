package parent

import "github.com/LucasBastino/app-sindicato/internal/common/page"

type tablePageData struct {
	Parents      []response
	TotalResults int
	EmptyState   page.EmptyState

	MemberID     int

	PageContext  page.PageContext
}

type pageData struct {
	Parent      response
	MemberID    int
	PageContext page.PageContext
	Errors      map[string]string
}
