package user

import "github.com/LucasBastino/app-sindicato/internal/common/page"

type PageData struct {
	User        Response
	PageContext page.PageContext
	Errors      map[string]string
}

type TablePageData struct {
	Users       []TableResponse
	PageContext page.PageContext
	Errors      map[string]string
}
