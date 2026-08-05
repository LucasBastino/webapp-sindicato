package company

import (
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/page"
)

type Company struct {
	ID     			 int       `db:"id_company"`

	Name             string    `db:"name"`
	CompanyNumber *string    `db:"company_number"`
	
	Address          string    `db:"address"`
	Cuit             *string    `db:"cuit"`
	District         string    `db:"district"`
	PostalCode       string    `db:"postal_code"`

	Phone            string    `db:"phone"`
	Contact          string    `db:"contact"`

	Observations     string    `db:"observations"`

	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	DeletedAt    	 *time.Time `db:"deleted_at"`
}

type companyFilters struct{
	searchKey string
	statuses page.StatusFilters
}

type RecentCompany struct {
	ID          int       `db:"id_company"`
	Name        string    `db:"name"`
	CreatedAt   time.Time `db:"created_at"`
	MemberCount int       `db:"member_count"`
}









