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

type PermissionsModalData struct {
	ID            int
	Username      string
	Admin         bool
	ResourceRoles map[string]any
	Errors        map[string]string
}

func (d PermissionsModalData) Role(resource string) string {
	if d.ResourceRoles == nil {
		return ""
	}
	v, _ := d.ResourceRoles[resource].(string)
	return v
}

type PasswordModalData struct {
	ID                     int
	RequireCurrentPassword bool
	CurrentPassword        string
	Password               string
	ConfirmPassword        string
	Errors                 map[string]string
}
