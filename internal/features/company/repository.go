package company

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type CompanyRepository struct{
	db *sqlx.DB
}

func NewCompanyRepository(db *sqlx.DB) *CompanyRepository{
	return &CompanyRepository{db: db}
}


func (r *CompanyRepository) FindByID(ctx context.Context, id int) (*Company, error) {
	query := "SELECT * FROM companies WHERE id_company = ?";
	var company Company
	err := r.db.GetContext(ctx, &company, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch company by id: %w", err)
	}
	return &company, nil
}

func (r *CompanyRepository) Search(ctx context.Context, filters companyFilters, offset int) ([]Company, error) {
	args := []any{}
	baseQuery := "SELECT * FROM companies"
	query, args := buildCompanyFilters(baseQuery, filters)
	query += " ORDER BY Name ASC LIMIT 15 OFFSET ?"
	args = append(args, offset)
	
	var companies []Company
	err := r.db.SelectContext(ctx, &companies, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search companies: %w", err)
	}
	return companies, nil
}

func (r *CompanyRepository) SearchForSelect(ctx context.Context, searchKey string) ([]Company, error){
	like := "%" + strings.TrimSpace(searchKey) + "%"

	query := `
	SELECT
		id_company,
		name
	FROM companies 
	WHERE ( 
		name LIKE ?
		OR company_number LIKE ?
	)
	AND deleted_at IS NULL 
	ORDER BY name ASC
	`

	var companies []Company
	err := r.db.SelectContext(ctx, &companies, query, like, like)
	if err!=nil{
		return nil, fmt.Errorf("failed to fetch companies for select: %w", err)
	}
	return companies, nil
}

func (r *CompanyRepository) ListActiveIDs(ctx context.Context) ([]int, error) {
	query := "SELECT id_company FROM companies WHERE is_deleted IS NULL "
	var ids []int
	err := r.db.SelectContext(ctx, &ids, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search active companies ids: %w", err)
	}
	return ids, nil
}

func (r *CompanyRepository) Count(ctx context.Context, filters companyFilters) (int, error) {
	args := []any{}
	baseQuery := "SELECT COUNT(*) FROM companies"
	query, args := buildCompanyFilters(baseQuery, filters)
	var totalRows int
	err := r.db.GetContext(ctx, &totalRows, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get total rows while searching companies: %w", err)
	}
	return totalRows, nil
}

func (r *CompanyRepository) FindRecent(ctx context.Context, limit int) ([]RecentCompany, error) {
	query := `
	SELECT
		C.id_company,
		C.name,
		C.created_at,
		(
			SELECT COUNT(*)
			FROM members M
			WHERE M.id_company = C.id_company
				AND M.deleted_at IS NULL
		) AS member_count
	FROM companies C
	WHERE C.deleted_at IS NULL
	ORDER BY C.created_at DESC
	LIMIT ?
	`
	var companies []RecentCompany
	err := r.db.SelectContext(ctx, &companies, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent companies: %w", err)
	}
	return companies, nil
}


func (r *CompanyRepository) Insert(ctx context.Context, tx *sqlx.Tx, company Company) (int, error) {

	query := `
	INSERT INTO companies (
		name,
		company_number,
		address,
		cuit,
		district,
		postal_code,
		phone,
		contact,
		observations
	) VALUES (
		:name,
		:company_number,
		:address,
		:cuit,
		:district,
		:postal_code,
		:phone,
		:contact,
		:observations
	)`
	
	// res, err := r.db.ExecContext(ctx, query, e.Name,
	// 	e.CompanyNumber,	e.Address, e.Cuit, e.District,
	// 	e.PostalCode, e.Phone, e.Contact, e.Observations)
	res, err := tx.NamedExecContext(ctx, query, company)
	if err != nil {
		return 0, fmt.Errorf("failed to insert company: %w", err)
	}
	id, err := res.LastInsertId()
	if err!=nil{
		return 0, fmt.Errorf("failed to get last insert id of companies: %w", err)
	}
	
	return int(id), nil
}

func (r *CompanyRepository) Update(ctx context.Context, id int, company Company) error {

	company.ID = id

	query := `
	UPDATE companies
	SET
		name = :name,
		company_number = :company_number,
		address = :address,
		cuit = :cuit,
		district = :district,
		postal_code = :postal_code,
		phone = :phone,
		contact = :contact,
		observations = :observations
	WHERE
		id_company = :id_company
		AND deleted_at IS NULL
	`
	
	_, err := r.db.NamedExecContext(ctx, query, company)
	if err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}
	return nil
}

func (r *CompanyRepository) SoftDelete(ctx context.Context, id int) (int, error){
	query := "UPDATE companies SET deleted_at = NOW() WHERE id_company = ?"
	res, err := r.db.ExecContext(ctx, query, id)
	if err!=nil{
		return 0, fmt.Errorf("failed to soft delete company: %w", err)
	}
	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while soft deleting company: %w", err)
	}
	return int(rows), nil
}

func (r *CompanyRepository) Restore(ctx context.Context, tx *sqlx.Tx, id int) (int, error){
	query := "UPDATE companies SET deleted_at = NULL WHERE id_company = ?"
	
	res, err := tx.ExecContext(ctx, query, id)
	if err!=nil{
		return 0, fmt.Errorf("failed to restore company: %w", err)
	}
	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while restoring company: %w", err)
	}
	return int(rows), nil
}

func (r *CompanyRepository) HardDelete(ctx context.Context, id int) error {
	query := `DELETE FROM companies WHERE id_company = ? AND deleted_at IS NOT NULL'`;
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete company: %w", err)
	}
	return nil
}


func (r *CompanyRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error){
	return r.db.BeginTxx(ctx, nil)
}



// validation functions


// func (r *CompanyRepository) ValidateCompanyInDB(ctx context.Context, id int, inputs Request) map[string]string {
// 	errorMap := map[string]string{}
// 	if err := r.validateCompanyNumber(ctx, id, inputs.CompanyNumber); err != "" {
//     	errorMap["companyNumber"] = err
// 	}
// 	if err := r.validateCuit(ctx, id, inputs.Cuit); err != "" {
//     	errorMap["cuit"] = err
// 	}
// 	return errorMap
// }


// func (r *CompanyRepository) validateCompanyNumber(ctx context.Context, id int, companyNumber string) string{
// 	var idDB int
// 	query := "SELECT id_company FROM companies WHERE company_number = ?";
// 	err := r.db.GetContext(ctx, &idDB, query, companyNumber)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return ""
// 		}
// 		return "Ocurrió un error con el campo."
// 	}
// 	if id != idDB{
// 		return "El número de empresa ingresado ya existe en el sistema."
// 	}
// 	return ""
// }

// func (r *CompanyRepository) validateCuit(ctx context.Context, id int, cuit string) string{
// 	var idDB int
// 	query := "SELECT id_company FROM companies WHERE cuit = ?";
// 	err := r.db.GetContext(ctx, &idDB, query, cuit)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return ""
// 		}
// 		return "Ocurrió un error con el campo."
// 	}
// 	if id != idDB{
// 		return "El CUIT ingresado ya existe en el sistema."
// 	}
// 	return ""
// }

// func (r *CompanyRepository) Exists(ctx context.Context, id int) (bool, error){
// 	query := "SELECT EXISTS (SELECT 1 FROM companies WHERE id_company = ?)"
// 	var exists bool
// 	err := r.db.GetContext(ctx, &exists, query, id)
// 	if err!=nil{
// 		return false, fmt.Errorf("failed to check if company exists: %w", err)
// 	}
// 	return exists, nil
// }