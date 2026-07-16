package userauthinfo

type UserAuthInfo struct {
	UserID int

	Admin         bool
	ResourceRoles map[string]any
}