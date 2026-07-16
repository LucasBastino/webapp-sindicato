package parent

import (
	"time"
)

type Parent struct {
	ID		     int       `db:"id_parent"`
	MemberID     int       `db:"id_member"`

	Name         string    `db:"name"`
	LastName     string    `db:"last_name"`
	Relationship string    `db:"relationship"`
	Birthday     time.Time `db:"birthday"`
	Gender       string    `db:"gender"`
	Cuil         *string   `db:"cuil"`

	Observations string		`db:"observations"`
	
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}




