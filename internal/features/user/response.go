package user

type Response struct {
	Username string

	Admin         bool
	ResourceRoles map[string]any
}

func (r Response) Role(resource string) string {
	if r.ResourceRoles == nil {
		return ""
	}
	v, _ := r.ResourceRoles[resource].(string)
	return v
}

type permissionsResponse struct {
	Admin         bool
	ResourceRoles map[string]any
}

type TableResponse struct {
	ID        int
	Username  string
	Admin     bool
	CreatedAt string
}
