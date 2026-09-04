package company

import (
	"strings"

	v "github.com/LucasBastino/webapp-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type request struct {
	Name                string `form:"name"`
	CompanyNumber    	string `form:"company-number"`

	Address             string `form:"address"`
	Cuit                string `form:"cuit"`
	District            string `form:"district"`
	PostalCode          string `form:"postal-code"`

	Phone               string `form:"phone"`
	Contact             string `form:"contact"`
	
	Observations        string `form:"observations"`
}


func (req *request) trim() {
	req.Name = strings.TrimSpace(req.Name)
	req.CompanyNumber = strings.TrimSpace(req.CompanyNumber)
	req.Address = strings.TrimSpace(req.Address)
	req.Cuit = strings.TrimSpace(req.Cuit)
	req.District = strings.TrimSpace(req.District)
	req.PostalCode = strings.TrimSpace(req.PostalCode)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Contact = strings.TrimSpace(req.Contact)
	req.Observations = strings.TrimSpace(req.Observations)
}

func (req request) validate() map[string]string {
	errorMap := map[string]string{}

	if err := v.ValidateCompanyName(req.Name); err != "" {
		errorMap["name"] = err
	}
	if err := v.ValidateCompanyNumber(req.CompanyNumber); err != "" {
		errorMap["companyNumber"] = err
	}
	if err := v.ValidateAddress(req.Address); err != "" {
		errorMap["address"] = err
	}
	if err := v.ValidateCuilCuit(req.Cuit); err != "" {
		errorMap["cuit"] = err
	}
	if err := v.ValidateDistrict(req.District); err != "" {
		errorMap["district"] = err
	}
	if err := v.ValidatePostalCode(req.PostalCode); err != "" {
		errorMap["postalCode"] = err
	}
	if err := v.ValidatePhone(req.Phone); err != "" {
		errorMap["phone"] = err
	}
	if err := v.ValidateContact(req.Contact); err != "" {
		errorMap["contact"] = err
	}
	if err := v.ValidateObservations(req.Observations); err != "" {
		errorMap["observations"] = err
	}

	return errorMap
}

