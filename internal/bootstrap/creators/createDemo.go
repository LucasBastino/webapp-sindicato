package creators

import "github.com/jmoiron/sqlx"

func CreateDemo(db *sqlx.DB) error {
	err := CreateCompanies(db)
	if err!=nil{
		return err
	}

	err = CreateMembers(db)
	if err!=nil{
		return err
	}

	err = CreateParents(db)
	if err!=nil{
		return err
	}
	
	err = CreatePayments(db)
	if err!=nil{
		return err
	}

	return nil
}