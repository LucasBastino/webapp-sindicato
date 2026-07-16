package user

type Response struct {
	Username string
	Password string

	Admin         bool
	ResourceRoles map[string]any
}

type passwordResponse struct {
	CurrentPassword string
	Password        string
	ConfirmPassword string
}

type permissionsResponse struct {
	Admin         bool
	ResourceRoles map[string]any
}

type TableResponse struct {
	Username  string
	Admin     bool
	CreatedAt string
}
