package member

import (
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/page"
)

type Member struct {
	ID      	  int       `db:"id_member"`
	CompanyID	  int       `db:"id_company"`

	Name          string    `db:"name"`
	LastName      string    `db:"last_name"`
	Dni           string    `db:"dni"`
	Birthday      time.Time `db:"birthday"`
	Gender        string    `db:"gender"`
	MaritalStatus string    `db:"marital_status"`

	Phone         string    `db:"phone"`
	Email         string    `db:"email"`

	Address       string    `db:"address"`
	PostalCode    string    `db:"postal_code"`
	District      string    `db:"district"`

	MemberNumber  string    `db:"member_number"`
	Cuil          *string    `db:"cuil"`
	Category      string    `db:"category"`
	EntryDate     time.Time `db:"entry_date"`

	Observations  string    `db:"observations"`

	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
	
	CompanyName     string     `db:"company_name"`
	CompanyDeletedAt *time.Time `db:"company_deleted_at"`
}

type memberFilters struct{
	searchKey string
	statuses page.StatusFilters
	companyID *int
}




