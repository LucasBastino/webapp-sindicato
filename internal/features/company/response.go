package company

type response struct {
	ID            int
	Name          string
	CompanyNumber string

	Address    string
	Cuit       string
	District   string
	PostalCode string

	Phone   string
	Contact string

	Observations string

	CreatedAt string
	UpdatedAt string

	IsDeleted bool
}

type tableResponse struct {
	ID            int
	Name          string
	CompanyNumber string
	Address       string
	IsDeleted     bool
}

type optionsResponse struct {
	ID   int
	Name string
}
