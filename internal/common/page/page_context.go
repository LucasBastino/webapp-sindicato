package page

import (
	"fmt"
	"net/url"

	userAuthInfo "github.com/LucasBastino/webapp-sindicato/internal/features/user/authinfo"
)

type PageContext struct {
	Mode          string
	UserAuthInfo  userAuthInfo.UserAuthInfo
	SearchKey     string
	Pagination    Pagination
	StatusFilters StatusFilters
	ActiveSection string // "dashboard" | "members" | "companies" | "users" | "reports"
}

func (p PageContext) TableQuery(page int) string {
	values := url.Values{}
	if page > 0 {
		values.Set("page", fmt.Sprintf("%d", page))
	}
	values.Set("show-active", fmt.Sprintf("%t", p.StatusFilters.ShowActive))
	values.Set("show-inactive", fmt.Sprintf("%t", p.StatusFilters.ShowInactive))
	values.Set("show-deleted", fmt.Sprintf("%t", p.StatusFilters.ShowDeleted))
	if p.SearchKey != "" {
		values.Set("search-key", p.SearchKey)
	}
	return values.Encode()
}

type StatusFilters struct {
	ShowActive   bool
	ShowInactive bool
	ShowDeleted  bool
}


// type ResourceRoles struct {
// 	Member     string
// 	Company string
// 	Parent     string
// 	Payment    string
// }