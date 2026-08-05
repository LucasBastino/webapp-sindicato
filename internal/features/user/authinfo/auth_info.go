package userauthinfo

type UserAuthInfo struct {
	UserID   int
	Username string

	Admin         bool
	ResourceRoles map[string]any
}

func (u UserAuthInfo) roleLevel(resource string) int {
	if u.Admin {
		return 3
	}
	if u.ResourceRoles == nil {
		return 0
	}
	role, ok := u.ResourceRoles[resource].(string)
	if !ok {
		return 0
	}
	switch role {
	case "editor":
		return 2
	case "viewer":
		return 1
	default:
		return 0
	}
}

// CanView reports whether the user can view the resource (Admin, viewer, or editor).
func (u UserAuthInfo) CanView(resource string) bool {
	return u.roleLevel(resource) >= 1
}

// CanEdit reports whether the user can edit the resource (Admin or editor).
func (u UserAuthInfo) CanEdit(resource string) bool {
	return u.roleLevel(resource) >= 2
}
