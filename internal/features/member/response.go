package member

type response struct {
	ID       int
	Name     string
	LastName string

	Dni           string
	Birthday      string
	Gender        string
	MaritalStatus string

	Phone string
	Email string

	Address    string
	PostalCode string
	District   string

	MemberNumber string
	Cuil         string
	CompanyID    int
	Category     string
	EntryDate    string

	Observations string

	CreatedAt string
	UpdatedAt string

	CompanyName string
}

type tableResponse struct {
	ID           int
	Name         string
	LastName     string
	Dni          string
	MemberNumber string
	CompanyName  string
	IsDeleted    bool
	IsInactive   bool
}
