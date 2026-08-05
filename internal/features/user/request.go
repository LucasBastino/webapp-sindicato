package user

import (
	"maps"

	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type Request struct {
	Username        string `form:"username"`

	passwordRequest
	permissionsRequest
}

type passwordRequest struct{
	CurrentPassword string `form:"current_password"`
	Password 		string `form:"password"`
	ConfirmPassword string `form:"confirm_password"`
}

type permissionsRequest struct {
	Admin   string `form:"admin"`
	Member  string `form:"member"`
	Company string `form:"company"`
}


func (req Request) Validate() map[string]string{
	errorMap := map[string]string{}

	if err := v.ValidateUsername(req.Username); err != "" {
		errorMap["username"] = err
	}

	maps.Copy(errorMap, req.passwordRequest.Validate(false))
	
	maps.Copy(errorMap, req.permissionsRequest.Validate())


	return errorMap
}

func (req passwordRequest) Validate(requireCurrent bool) map[string]string{
	errorMap := map[string]string{}

	if requireCurrent {
		if req.CurrentPassword == "" {
			errorMap["current_password"] = "Campo requerido."
		}
	}

	if err := v.ValidatePassword(req.Password); err != "" {
		errorMap["password"] = err
	}
	if req.ConfirmPassword == "" {
		errorMap["confirm_password"] = "Campo requerido."
		return errorMap

	} else if req.Password != req.ConfirmPassword {
		errorMap["confirm_password"] = "Las contraseñas no coinciden."
		errorMap["password"] = "Las contraseñas no coinciden."
	}

	return errorMap
}

func (req permissionsRequest) Validate() map[string]string{
	errorMap := map[string]string{}

		if err := v.ValidateAdmin(req.Admin); err != "" {
		errorMap["admin"] = err
	}

	if err := v.ValidatePermissions(req.Member); err != "" {
		errorMap["member"] = err
	}

	if err := v.ValidatePermissions(req.Company); err != "" {
		errorMap["company"] = err
	}

	return errorMap
}

	
