package parent

import (
	"strings"

	v "github.com/LucasBastino/webapp-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type request struct {
	Name         string `form:"name"`
	LastName     string `form:"last-name"`
	Relationship string `form:"relationship"`
	Birthday     string `form:"birthday"`
	Gender       string `form:"gender"`
	Cuil         string `form:"cuil"`
	Observations string `form:"observations"`
}

func (req *request) trim() {
	req.Name = strings.TrimSpace(req.Name)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Relationship = strings.TrimSpace(req.Relationship)
	req.Birthday = strings.TrimSpace(req.Birthday)
	req.Gender = strings.TrimSpace(req.Gender)
	req.Cuil = strings.TrimSpace(req.Cuil)
	req.Observations = strings.TrimSpace(req.Observations)
}

// no chequeo memberID porque el parent no puede cambiar de member
func (req request) validate() map[string]string {

	errorMap := map[string]string{}

	if err := v.ValidateNameOrLastName(req.Name); err != "" {
		errorMap["name"] = err
	}
	if err := v.ValidateNameOrLastName(req.LastName); err != "" {
		errorMap["lastName"] = err
	}
	if err := v.ValidateOptionalRelationship(req.Relationship); err != "" {
		errorMap["relationship"] = err
	}
	if err := v.ValidateOptionalBirthday(req.Birthday); err != "" {
		errorMap["birthday"] = err
	}
	if err := v.ValidateOptionalGender(req.Gender); err != "" {
		errorMap["gender"] = err
	}
	if err := v.ValidateCuilCuit(req.Cuil); err != "" {
		errorMap["cuil"] = err
	}

	return errorMap
}

