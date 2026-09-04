package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository{
	return &UserRepository{db: db}
}

// func (r *UserRepository) Exists(ctx context.Context, username string) (bool, error){
// 	query := "SELECT EXISTS (SELECT 1 FROM users WHERE username = ?)"
// 	var exists bool
// 	err := r.db.GetContext(ctx, &exists, query, username)
// 	if err!=nil{
// 		return false, fmt.Errorf("failed to check if user exists: %w", err)
// 	}
// 	return exists, nil
// }


func (r *UserRepository) FindByID(ctx context.Context, id int) (*User, error) {
	query := "SELECT id_user, username, password_hash, admin, resource_roles, created_at FROM users WHERE id_user = ?"
	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}
	// Convierto los permisos JSON en un map[string]any
	err = json.Unmarshal(user.ResourceRolesJSON, &user.ResourceRoles)
	if err!=nil{
		return nil, fmt.Errorf("failed to unmarshal resource roles while fetching user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE username = ?"
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil
        }
		return nil, fmt.Errorf("failed to fetch user by username: %w", err)
	}
	// Convierto los permisos JSON en un map[string]any
	err = json.Unmarshal(user.ResourceRolesJSON, &user.ResourceRoles)
	if err!=nil{
		return nil, fmt.Errorf("failed to unmarshal resource roles while fetching user by username: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]User, error) {
    query := "SELECT id_user, username, admin, resource_roles, created_at FROM users"
    var users []User
    err := r.db.SelectContext(ctx, &users, query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch users: %w", err)
    }
    return users, nil
}


func (r *UserRepository) Insert(ctx context.Context, tx *sqlx.Tx, user User) (int, error) {
	query := "INSERT INTO users (username, password_hash, admin, resource_roles) VALUES (:username, :password_hash, :admin, :resource_roles)"
	res, err := tx.NamedExecContext(ctx, query, user)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	id, err := res.LastInsertId()
	if err!=nil{
		return 0, fmt.Errorf("failed to get last insert id while inserting user: %w", err)
	}
	return int(id), nil
}

func (r *UserRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}


// func (r *UserRepository) GetRoleByID(ctx context.Context, userID int) (ResourceRoles, error){
// 	query := "SELECT is_admin, can_write, can_delete FROM users WHERE id_user = ?"
// 	var permissions ResourceRoles
// 	err := r.db.GetContext(ctx, &permissions, query, userID)
// 	if err!=nil{
// 		 if errors.Is(err, sql.ErrNoRows) {
//             return ResourceRoles{}, apperrors.NewNotFoundError(fmt.Errorf("user not found: %w", err), "")
//         }
// 		return ResourceRoles{}, fmt.Errorf("failed to fetch user ResourceRoles %w", err)
// 	}
// 	return permissions, nil

// }



func (r *UserRepository) UpdatePermissions(ctx context.Context, id int, admin bool, resourceRoles map[string]any) error {
	rolesJSON, err := json.Marshal(resourceRoles)
	if err != nil {
		return fmt.Errorf("failed to marshal resource roles: %w", err)
	}
	query := "UPDATE users SET admin = ?, resource_roles = ? WHERE id_user = ?"
	_, err = r.db.ExecContext(ctx, query, admin, rolesJSON, id)
	if err != nil {
		return fmt.Errorf("failed to update user permissions: %w", err)
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id int, hash string) (int, error) {
    query := "UPDATE users SET password_hash = ? WHERE id_user = ?"
    res, err := r.db.ExecContext(ctx, query, hash, id)
    if err != nil {
        return 0, fmt.Errorf("failed to update password: %w", err)
    }
    rows, err := res.RowsAffected()
	if err!=nil{
		return 0, fmt.Errorf("failed to get rows affected while updating user password: %w", err)
	}
	return int(rows), nil
}

/* func (r *UserRepository) Getpassword_hashFromUsername(username string) (string, error){
	row := r.db.QueryRow("SELECT password_hash FROM UserTable WHERE Username = ?", username)
	var password_hashStr string
	err := row.Scan(&password_hashStr)
	if err!=nil{
		return "", err
	}
	return password_hashStr, nil
} */

func (r *UserRepository) CountAdmins(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM users WHERE admin = true"
	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin users: %w", err)
	}
	return count, nil
}

func (r *UserRepository) HardDelete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id_user = ?";
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to hard delete user: %w", err)
	}	
	return nil
}
