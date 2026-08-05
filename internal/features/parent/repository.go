package parent

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ParentRepository struct{
	db *sqlx.DB
}

func NewParentRepository(db *sqlx.DB) *ParentRepository{
	return &ParentRepository{db :db}
}


func (r *ParentRepository) FindByID(ctx context.Context, id int) (*Parent, error) {
	query := "SELECT * FROM parents WHERE id_parent = ?"

	var parent Parent

	err := r.db.GetContext(ctx, &parent, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch parent by id: %w", err)
	}

	return &parent, nil
}

func (r *ParentRepository) FindAll(ctx context.Context, memberID int) ([]Parent, error) {
	query := "SELECT * FROM parents WHERE id_member = ?"
	var parents []Parent
	err := r.db.SelectContext(ctx, &parents, query, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch parents: %w", err)
	}
	return parents, nil
}

func (r *ParentRepository) Count(ctx context.Context, memberID int) (int, error) {
	query := "SELECT COUNT(*) FROM parents WHERE id_member = ?"
	var totalRows int
	err := r.db.GetContext(ctx, &totalRows, query, memberID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total rows while fetching parents: %w", err)
	}
	return totalRows, nil
}


func (r *ParentRepository) Insert(ctx context.Context, tx *sqlx.Tx, parent Parent) (int, error) {
	query := `
		INSERT INTO parents (
		id_member,
		name,
		last_name,
		relationship,
		birthday,
		gender,
		cuil,
		observations
		)
		VALUES (
		:id_member,
		:name,
		:last_name,
		:relationship,
		:birthday,
		:gender,
		:cuil,
		:observations
		)`;
	res, err := tx.NamedExecContext(ctx, query, parent)
	if err != nil {
		return 0, fmt.Errorf("failed to insert parent: %w", err)
	}
	id, err := res.LastInsertId()
	if err!=nil{
		return 0, fmt.Errorf("failed to get last insert id of parent: %w", err)
	}
	return int(id), nil
}

func (r *ParentRepository) Update(ctx context.Context, id int, parent Parent) error {
	parent.ID = id
	query := `
	UPDATE parents
	SET
    name = :name,
    last_name = :last_name,
    relationship = :relationship,
    birthday = :birthday,
    gender = :gender,
    cuil = :cuil,
    observations = :observations
	WHERE id_parent = :id_parent`
	_, err := r.db.NamedExecContext(ctx, query, parent)
	if err != nil {
		return fmt.Errorf("failed to update parent: %w", err)
	}
	return nil
}

func (r *ParentRepository) HardDelete(ctx context.Context, id int) error {
	query := "DELETE FROM parents WHERE id_parent = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete parent: %w", err)
	}
	return nil
}


func (r *ParentRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error){
	return r.db.BeginTxx(ctx, nil)
}

// func (r *ParentRepository) HardDeleteAllParentsByMember(ctx context.Context, memberID int) error{
// 	query := "DELETE FROM parent WHERE member_id = ?"
// 	// doesn't need verify the number of parents deleted because they aren't important models
// 	// and the member couldn't have any parents
// 	_, err := r.db.ExecContext(ctx, query)
// 	if err!= nil{
// 		return fmt.Errorf("failed to hard delete parents: %w", err)
// 	}
// 	return nil
// }

// func (r *ParentRepository) ValidateInDB(ctx context.Context, parentID, memberID int, inputs Request) map[string]string{
// 	errorMap := map[string]string{}
// 	if err := r.validateMemberID(ctx, memberID); err!=""{
// 		errorMap["memberID"] = err
// 	}
// 	if err := r.validateCuil(ctx, parentID, inputs.Cuil); err!=""{
// 		errorMap["cuil"] = err
// 	}
// 	return errorMap
// }

// func (r *ParentRepository) validateMemberID(ctx context.Context, memberID int) string{

// 	var dummy int
// 	query := "SELECT 1 FROM member WHERE id_member = ?";
// 	err := r.db.GetContext(ctx, &dummy, query, memberID)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return "Afiliado no existente."
// 		}
// 		return "Ocurrió un error con el campo."

// 	}
// 	return ""
// }

// func (r *ParentRepository) validateCuil(ctx context.Context, id int, cuil string) string{
// 	query := "SELECT id_parent FROM parent WHERE cuil = ?"
// 	var idDB int
// 	err := r.db.GetContext(ctx, &idDB, query, cuil)
// 	if err!=nil{
// 		if errors.Is(err, sql.ErrNoRows){
// 			return ""
// 		}
// 		return "Ocurrió un error con el campo."
// 	}
// 	if id != idDB{
// 		return "El CUIL ya fue ingresado."
// 	}
// 	return ""
// }






// func (r *ParentRepository) Exists(ctx context.Context, id int) (bool, error){
// 	query := "SELECT EXISTS (SELECT 1 FROM parent WHERE id_parent = ?)"
// 	var exists bool
// 	err := r.db.GetContext(ctx, &exists, query, id)
// 	if err!=nil{
// 		return false, fmt.Errorf("failed to check if parent exists: %w", err)
// 	}
// 	return exists, nil
// }

