package page

// CompanyNav feeds the shared company sidecard (data + management nav).
// Active is one of: ficha | members | payments | plans
type CompanyNav struct {
	CompanyID      int
	CompanyName    string
	CompanyNumber  string
	Address        string
	Phone          string
	Active         string
	CanViewMembers bool
}

func NewCompanyNav(companyID int, companyName, active string, canViewMembers bool) CompanyNav {
	return CompanyNav{
		CompanyID:      companyID,
		CompanyName:    companyName,
		Active:         active,
		CanViewMembers: canViewMembers,
	}
}

func (n CompanyNav) WithDetails(number, address, phone string) CompanyNav {
	n.CompanyNumber = number
	n.Address = address
	n.Phone = phone
	return n
}
