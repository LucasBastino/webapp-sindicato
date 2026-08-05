package member

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MemberRepository struct{
	db *sqlx.DB
	
}

func NewMemberRepository(db *sqlx.DB) *MemberRepository{
	return &MemberRepository{db: db}
}


func (r *MemberRepository) FindByID(ctx context.Context, id int) (*Member, error) {
	query := `
	SELECT 
		M.*,
		C.name AS company_name,
		C.deleted_at AS company_deleted_at
	FROM members M
	INNER JOIN companies C
		ON M.id_company = C.id_company
	WHERE
		id_member = ?
	`
	var member Member
	err := r.db.GetContext(ctx, &member, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch member by id: %w", err)
	}
	return &member, nil
}

func (r *MemberRepository) Search(ctx context.Context, filters memberFilters, offset int) ([]Member, error) {
	args := []any{}
	baseQuery := `
	SELECT
		M.id_member,
		M.name,
		M.last_name,
		M.dni,
		M.member_number,
		M.deleted_at,
		C.name AS company_name,
		C.deleted_at AS company_deleted_at
	FROM members M
	INNER JOIN companies C
		ON M.id_company = C.id_company
	`
	
	query, args := buildMemberFilters(baseQuery, filters)

	// siempre dejar espacio antes
	query += ` ORDER BY last_name ASC LIMIT 15 OFFSET ?`
	args = append(args, offset)

	var members []Member
	err := r.db.SelectContext(ctx, &members, query, args...)
	if err != nil {
		return nil,  fmt.Errorf("failed to search members: %w", err)
	}
	return members, nil
}

func (r *MemberRepository) GetElectoralMemberList(ctx context.Context) ([]Member, error){
	var members []Member
	query := `
	SELECT
		M.name,
		M.last_name,
		M.dni,
		M.member_number,
		C.name AS company_name
	FROM members M
	INNER JOIN companies C
		ON M.id_company = C.id_company
	WHERE
		C.deleted_at IS NULL
		AND	M.deleted_at IS NULL
	ORDER BY M.last_name ASC
	`
	err := r.db.SelectContext(ctx, &members, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch electoral member list: %w", err)
	}
	return members, nil
}

func (r *MemberRepository) Count(ctx context.Context, filters memberFilters) (int, error) {
	args := []any{}
	baseQuery :="SELECT COUNT(*) FROM members M JOIN companies C ON M.id_company = C.id_company"
	query, args := buildMemberFilters(baseQuery, filters)
	var totalRows int
	err := r.db.GetContext(ctx, &totalRows, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get total rows while searching members: %w", err)
	}
	return totalRows, nil
}

func (r *MemberRepository) FindRecent(ctx context.Context, limit int) ([]Member, error) {
	query := `
	SELECT
		M.id_member,
		M.name,
		M.last_name,
		M.created_at,
		C.name AS company_name
	FROM members M
	INNER JOIN companies C
		ON M.id_company = C.id_company
	WHERE
		M.deleted_at IS NULL
		AND C.deleted_at IS NULL
	ORDER BY M.created_at DESC
	LIMIT ?
	`
	var members []Member
	err := r.db.SelectContext(ctx, &members, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent members: %w", err)
	}
	return members, nil
}


func (r *MemberRepository) Insert(ctx context.Context, tx *sqlx.Tx, member Member) (int, error) {
	query := 
	`INSERT INTO members (
		name,
		last_name,
		dni,
		birthday,
		gender,
		marital_status,
		phone,
		email,
		address,
		postal_code,
		district,
		member_number,
		cuil,
		id_company,
		category,
		entry_date,
		observations
	) VALUES (
		:name,
		:last_name,
		:dni,
		:birthday,
		:gender,
		:marital_status,
		:phone,
		:email,
		:address,
		:postal_code,
		:district,
		:member_number,
		:cuil,
		:id_company,
		:category,
		:entry_date,
		:observations
	)`
	res, err := tx.NamedExecContext(ctx, query, member)
	if err != nil {
		return 0, fmt.Errorf("failed to insert member: %w", err)
	}
	id, err := res.LastInsertId()
	if err!=nil{
		return 0, fmt.Errorf("failed to get last insert id of member: %w", err)
	}
	return int(id), nil
}

func (r *MemberRepository) Update(ctx context.Context, id int, member Member) error {

	member.ID = id

	query := `
	UPDATE members
	SET
		name = :name,
		last_name = :last_name,
		dni = :dni,
		birthday = :birthday,
		gender = :gender,
		marital_status = :marital_status,
		phone = :phone,
		email = :email,
		address = :address,
		postal_code = :postal_code,
		district = :district,
		member_number = :member_number,
		cuil = :cuil,
		id_company = :id_company,
		category = :category,
		entry_date = :entry_date,
		observations = :observations
	WHERE
		id_member = :id_member
		AND deleted_at IS NULL
	`

	_, err := r.db.NamedExecContext(ctx, query, member)
	if err != nil {
		return fmt.Errorf("failed to update member: %w", err)
	}
	return nil
}

func (r *MemberRepository) SoftDelete(ctx context.Context, id int) (int, error){
	query := "UPDATE members SET deleted_at = NOW() WHERE id_member = ?"
	res, err := r.db.ExecContext(ctx, query, id)
	if err!=nil{
		return 0, fmt.Errorf("failed to soft delete member: %w", err)
	}	
	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while soft deleting member: %w", err)
	}

	return int(rows), nil
}

func (r *MemberRepository) Restore(ctx context.Context, id int) (int, error){
	query := "UPDATE members SET deleted_at = NULL WHERE id_member = ?"
	res, err :=  r.db.ExecContext(ctx, query, id)
	if err!=nil{
		return 0, fmt.Errorf("failed to restore member: %w", err)
	}
	rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while restoring member: %w", err)
	}
	return int(rows), nil
}

func (r *MemberRepository) HardDelete(ctx context.Context, id int) error {
	query := "DELETE FROM members WHERE id_member = ? AND deleted_at IS NOT NULL";
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete member: %w", err)
	}	
	return nil
}


func (r *MemberRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error){
	return r.db.BeginTxx(ctx, nil)
}

// func (r *MemberRepository) FindByCompany(ctx context.Context, companyID, limit, offset int, searchKey string, includeInactive, includeDeleted bool) ([]Member, error) {
// 	args := []any{}
// 	args = append(args, companyID)
// 	baseQuery := `
// 		SELECT 
// 		M.name, 
// 		M.last_name, 
// 		M.dni 
// 		FROM member M
// 		INNER JOIN companies E ON M.id_company = C.id_company
// 		WHERE M.id_company = ?`
// 	query, args, err := buildMemberFilters(baseQuery, searchKey, includeInactive, includeDeleted)
// 	args = append(args, limit, offset)
// 	query += " LIMIT ? OFFSET ?"
// 	var members []Member
// 	err = r.db.SelectContext(ctx, &members, query, args...)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to search members by company: %w", err)
// 	}
// 	return members, nil
// }


// func (r *MemberRepository) ValidateInDB(ctx context.Context, idMember int, inputs Request) map[string]string{
// 	errorMap := map[string]string{}
// 	if inputs.CompanyID != ""{
// 		if err := r.validateCompanyID(ctx, inputs.CompanyID); err != "" {
// 			errorMap["companyID"] = err
// 		}
// 	}
// 	if err := r.validateMemberNumber(ctx, idMember, inputs.MemberNumber); err != "" {
// 		errorMap["memberNumber"] = err
// 	}
// 	if err := r.validateCuil(ctx, idMember, inputs.Cuil); err != "" {
// 		errorMap["cuil"] = err
// 	}

// 	return errorMap
// }

// func (r *MemberRepository) validateCompanyID(ctx context.Context, companyID string) string{

// 	var dummy int
// 	query := "SELECT 1 FROM companies WHERE id_company = ?";
// 	err := r.db.GetContext(ctx, &dummy, query, companyID)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return "Empresa no existente."
// 		}
// 		return "Ocurrió un error con el campo."

// 	}
// 	return ""
// }

// func (r *MemberRepository) validateMemberNumber(ctx context.Context, id int, memberNumber string) string{
// 	var idDB int
// 	query := "SELECT id_member FROM member WHERE member_number = ?";
// 	err := r.db.GetContext(ctx, &idDB, query, memberNumber)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return ""
// 		}
// 		return "Ocurrió un error con el campo."
// 	}
// 	if id != idDB{
// 		return "El número de afiliado ingresado ya existe en el sistema."
// 	}
// 	return ""
// }

// func (r *MemberRepository) validateCuil(ctx context.Context, id int, cuil string) string{
// 	var idDB int
// 	query := "SELECT id_member FROM member WHERE cuil = ?";
// 	err := r.db.GetContext(ctx, &idDB, query, cuil)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return ""
// 		}
// 		return "Ocurrió un error con el campo."
// 	}
// 	if id != idDB{
// 		return "Este CUIL ya ha sido ingresado."
// 	}
// 	return ""
// }













// func (r *MemberRepository) Exists(ctx context.Context, id int) (bool, error){
// 	query := "SELECT EXISTS (SELECT 1 FROM member WHERE id_member = ?)"
// 	var exists bool
// 	err := r.db.GetContext(ctx, &exists, query, id)
// 	if err!=nil{
// 		return false, fmt.Errorf("failed to check if member exists: %w", err)
// 	}
// 	return exists, nil
// }






