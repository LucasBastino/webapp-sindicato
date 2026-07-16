package member

import (
	"strings"

	v "github.com/LucasBastino/app-sindicato/internal/validation"
)

// recordar que los campos tiene que ser exportados para que el body parser los lea
type request struct {
	Name            string `form:"name"`
	LastName        string `form:"last-name"`
	Dni             string `form:"dni"`
	Birthday        string `form:"birthday"`
	Gender          string `form:"gender"`
	MaritalStatus   string `form:"marital-status"`

	Phone           string `form:"phone"`
	Email           string `form:"email"`

	Address         string `form:"address"`
	PostalCode      string `form:"postal-code"`
	District        string `form:"district"`

	MemberNumber    string `form:"member-number"`
	Cuil            string `form:"cuil"`
	CompanyID    string `form:"company-id"`
	Category        string `form:"category"`
	EntryDate       string `form:"entry-date"`

	Observations    string `form:"observations"`
	
	CompanyName  string `form:"company-name"`
}

func (req *request) trim() {
	req.Name = strings.TrimSpace(req.Name)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Dni = strings.TrimSpace(req.Dni)
	req.Birthday = strings.TrimSpace(req.Birthday)
	req.Gender = strings.TrimSpace(req.Gender)
	req.MaritalStatus = strings.TrimSpace(req.MaritalStatus)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Address = strings.TrimSpace(req.Address)
	req.PostalCode = strings.TrimSpace(req.PostalCode)
	req.District = strings.TrimSpace(req.District)
	req.MemberNumber = strings.TrimSpace(req.MemberNumber)
	req.Cuil = strings.TrimSpace(req.Cuil)
	req.CompanyID = strings.TrimSpace(req.CompanyID)
	req.Category = strings.TrimSpace(req.Category)
	req.EntryDate = strings.TrimSpace(req.EntryDate)
	req.Observations = strings.TrimSpace(req.Observations)
}

func (req request) validate() map[string]string {
	
	errorMap := map[string]string{}
	
	if err := v.ValidateNameOrLastName(req.Name); err != "" {
		errorMap["name"] = err
	}
	if err := v.ValidateNameOrLastName(req.LastName); err != "" {
		errorMap["lastName"] = err
	}
	if err := v.ValidateDni(req.Dni); err != "" {
		errorMap["dni"] = err
	}
	if err := v.ValidateBirthday(req.Birthday); err != "" {
		errorMap["birthday"] = err
	}
	if err := v.ValidateGender(req.Gender); err != "" {
		errorMap["gender"] = err
	}
	if err := v.ValidateMaritalStatus(req.MaritalStatus); err != "" {
		errorMap["maritalStatus"] = err
	}
	if err := v.ValidatePhone(req.Phone); err != "" {
		errorMap["phone"] = err
	}
	if err := v.ValidateEmail(req.Email); err != "" {
		errorMap["email"] = err
	}
	if err := v.ValidateAddress(req.Address); err != "" {
		errorMap["address"] = err
	}
	if err := v.ValidatePostalCode(req.PostalCode); err != "" {
		errorMap["postalCode"] = err
	}
	if err := v.ValidateDistrict(req.District); err != "" {
		errorMap["district"] = err
	}
	if err := v.ValidateMemberNumber(req.MemberNumber); err != "" {
		errorMap["memberNumber"] = err
	}
	if err := v.ValidateCuilCuit(req.Cuil); err != "" {
		errorMap["cuil"] = err
	}
	// chequeo companyID porque el afiliado puede cambiar de empresa
	if err := v.ValidateCompanyID(req.CompanyID); err != "" {
		errorMap["companyID"] = err
	}
	if err := v.ValidateCategory(req.Category); err != "" {
		errorMap["category"] = err
	}
	if err := v.ValidateEntryDate(req.EntryDate); err != "" {
		errorMap["entryDate"] = err
	}
	if err := v.ValidateObservations(req.Observations); err != "" {
		errorMap["observations"] = err
	}

	return errorMap
}

