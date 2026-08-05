package user

func ToModel(req Request) User {
	admin, roles := req.ToPermissions()

	return User{
		Username:      req.Username,
		Admin:         admin,
		ResourceRoles: roles,
	}
}

func ToResponseFromRequest(req Request) Response {
	admin, roles := req.ToPermissions()

	return Response{
		Username:        req.Username,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Admin:           admin,
		ResourceRoles:   roles,
	}
}

func ToPasswordResponseFromRequest(req passwordRequest) passwordResponse {
	return passwordResponse{
		CurrentPassword: req.CurrentPassword,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}
}

func ToPermissionsResponseFromRequest(req permissionsRequest) permissionsResponse {
	admin, roles := req.ToPermissions()
	return permissionsResponse{
		Admin:         admin,
		ResourceRoles: roles,
	}
}

func ToTableResponse(u User) TableResponse {
	return TableResponse{
		ID:        u.ID,
		Username:  u.Username,
		Admin:     u.Admin,
		CreatedAt: u.CreatedAt.Format("02/01/2006"),
	}
}

func ToTableResponses(users []User) []TableResponse {
	responses := make([]TableResponse, len(users))
	for i, user := range users {
		responses[i] = ToTableResponse(user)
	}
	return responses
}

func (req permissionsRequest) ToPermissions() (bool, map[string]any) {
	admin := req.Admin == "on"

	member := req.Member
	company := req.Company

	if admin {
		member = "editor"
		company = "editor"
	}

	roles := map[string]any{
		"member":  member,
		"company": company,
	}

	return admin, roles
}