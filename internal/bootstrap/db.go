package bootstrap

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func initDB(db *sqlx.DB) error {
	err := ensureTables(db)
	if err!=nil{
		return err
	}

	err = ensureIndexes(db)
	if err!=nil{
		return err
	}

	return nil
}


func ensureTables(db *sqlx.DB) error{
	err := ensureUsers(db)
	if err!=nil {
		return err
	}

	err = ensureRefreshTokens(db)
	if err != nil {
		return err
	}

	err = ensureCompanies(db)
	if err!=nil {
		return err
	}

	err = ensureMembers(db)
	if err!=nil {
		return err
	}

	err = ensureParents(db)
	if err!=nil {
		return err
	}

	err = ensurePayments(db)
	if err != nil {
		return err
	}

	err = ensurePaymentPlans(db)
	if err != nil {
		return err
	}

	err = ensureInstallments(db)
	if err != nil {
		return err
	}

	err = ensureIdempotencyKeys(db)
	if err!=nil {
		return err
	}
	
	hasUsers, err := hasUsers(db)
	if err!=nil {
		return err
	}
	if !hasUsers {
		err = ensureDefaultAdmin(db)
		if err != nil {
			return err
		}
	}
	
	return nil
}

func ensureIndexes(db *sqlx.DB) error {
	err := ensureCompanyIndexes(db)
	if err!=nil{
		return err
	}

	err = ensureMemberIndexes(db)
	if err!=nil{
		return err
	}

	err = ensureParentIndexes(db)
	if err!=nil{
		return err
	}

	err = ensurePaymentIndexes(db)
	if err!=nil{
		return err
	}

	err = ensurePaymentPlanIndexes(db)
	if err!=nil{
		return err
	}

	err = ensureIdempotencyKeysIndexes(db)
	if err!=nil{
		return err
	}

	err = ensureRefreshTokenIndexes(db)
	if err != nil {
		return err
	}

	return nil
}



func ensureCompanies(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS companies(
		id_company INT PRIMARY KEY AUTO_INCREMENT NOT NULL,

		name VARCHAR(150) NOT NULL,
		company_number VARCHAR(20),
		cuit VARCHAR(13),

		address VARCHAR(50),
		district VARCHAR(50),
		postal_code VARCHAR(8),

		phone VARCHAR(50),
		contact VARCHAR(200),

		observations VARCHAR(2000),

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP DEFAULT NULL,

		CONSTRAINT unique_company_number UNIQUE (company_number),
		CONSTRAINT unique_cuit UNIQUE (cuit)
	)`)

	if err!=nil{
		return fmt.Errorf("failed to create companies table: %w", err)
	}
	return nil
}

func ensureCompanyIndexes(db *sqlx.DB) error {
	// company_number ya tiene índice por constraint UNIQUE (unique_company_number)
	queries := []string{
	`
	CREATE INDEX idx_company_deleted_at
	ON companies(deleted_at)
	`,
	`
	CREATE INDEX idx_company_name
	ON companies(name)
	`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			// 1061 = duplicate key
			if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1061) {
				continue
			}
            return fmt.Errorf("failed to create companies indexes: %w", err)
        }
    }

    return nil
}


func ensureMembers(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS members(
		id_member INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		id_company INT,

		name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		dni VARCHAR(8),
		birthday DATE,
		gender VARCHAR(30) NOT NULL,
		marital_status VARCHAR(20) NOT NULL,

		phone VARCHAR(20),
		email VARCHAR(50),

		address VARCHAR(100),
		postal_code VARCHAR(8),
		district VARCHAR(50),

		member_number VARCHAR(20),
		cuil VARCHAR(13),
		category VARCHAR(100),
		entry_date DATE,

		observations VARCHAR(2000),

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP,

		CONSTRAINT unique_dni_gender UNIQUE (dni, gender),
		CONSTRAINT unique_member_number UNIQUE (member_number),
		CONSTRAINT unique_cuil UNIQUE (cuil),

		FOREIGN KEY (id_company) REFERENCES companies(id_company) ON DELETE SET NULL
	)`)
	// id_company va sin NOT NULL, por si borras la empresa
	if err!=nil{
		return fmt.Errorf("failed to create members table: %w", err)
	}
	_, err = db.Exec(`UPDATE members SET member_number = NULL WHERE member_number = ''`)
	if err != nil {
		return fmt.Errorf("failed to normalize empty member_number to NULL: %w", err)
	}
	return nil
}

func ensureMemberIndexes(db *sqlx.DB) error {
	// member_number ya tiene índice por constraint UNIQUE (unique_member_number)
	queries := []string{
	`
	CREATE INDEX idx_member_id_company_deleted_at
	ON members(id_company, deleted_at)
	`,
	`
	CREATE INDEX idx_member_deleted_at
	ON members(deleted_at)
	`,
	`
	CREATE INDEX idx_member_name
	ON members(name)
	`,
	`
	CREATE INDEX idx_member_last_name
	ON members(last_name)
	`,
	// dni no tiene UNIQUE simple; el UNIQUE es compuesto (dni, gender)
	`
	CREATE INDEX idx_member_dni
	ON members(dni)
	`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			// 1061 = duplicate key
			if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1061) {
				continue
			}
            return fmt.Errorf("failed to create members indexes: %w", err)
        }
    }

    return nil
}


func ensureParents(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS parents(
		id_parent INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		id_member INT NOT NULL,

		name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		relationship VARCHAR(20) NOT NULL,
		birthday DATE,
		gender VARCHAR(50),
		cuil VARCHAR(13),

		observations VARCHAR(2000),

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT unique_cuil UNIQUE (cuil),

		FOREIGN KEY (id_member) REFERENCES members(id_member) ON DELETE CASCADE
	)`)
	// aca si id_member va NOT NULL porque si borro el member que se borren los parent tambien
	if err!=nil{
		return fmt.Errorf("failed to create parents table: %w", err)
	}
	return nil
}

func ensureParentIndexes(db *sqlx.DB) error {
	query := 
	`
	CREATE INDEX idx_parent_id_member
	ON parents(id_member)
	`

	_, err := db.Exec(query)
	if err != nil {
		var mysqlErr *mysql.MySQLError
			// 1061 = duplicate key
			if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1061) {
				return nil
			}
		return fmt.Errorf("failed to create parents indexes: %w", err)
	}

    return nil
}


func ensurePayments(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS payments(
		id_payment INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		id_company INT NOT NULL,

		month INT UNSIGNED NOT NULL,
		year INT UNSIGNED NOT NULL,
		due_date DATE NOT NULL,
		amount DECIMAL(10,2),
		paid_at DATE DEFAULT null,
		is_in_payment_plan BOOLEAN NOT NULL DEFAULT false ,

		observations VARCHAR(2000),

		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT unique_id_company_month_year UNIQUE (id_company, month, year),

		FOREIGN KEY (id_company) REFERENCES companies(id_company) ON DELETE CASCADE
	)`)
	if err!=nil{
		return fmt.Errorf("failed to create payments table: %w", err)
	}
	return nil
}

func ensurePaymentIndexes(db *sqlx.DB) error {
	query := 
	`
	CREATE INDEX idx_payment_year_month
	ON payments(year, month)
	`

	_, err := db.Exec(query)
	if err != nil {
		var mysqlErr *mysql.MySQLError
			// 1061 = duplicate key
			if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1061) {
				return nil
			}
		return fmt.Errorf("failed to create payments indexes: %w", err)
	}

    return nil
}


func ensurePaymentPlans(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS payment_plans(
		id_payment_plan INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		id_company INT NOT NULL,

		payments_in_plan VARCHAR(200),
		amount DECIMAL NOT NULL,
		number_of_installments INT,
		status VARCHAR(20) NOT NULL DEFAULT 'pending',
		first_due_date DATE NOT NULL,
		last_due_date DATE NOT NULL,

		observations VARCHAR(2000),

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		FOREIGN KEY (id_company) REFERENCES companies(id_company) ON DELETE CASCADE
	)`)
	if err != nil {
		return fmt.Errorf("failed to create payment plans table: %w", err)
	}

	// Existing DBs created before status existed.
	_, err = db.Exec(`ALTER TABLE payment_plans ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending'`)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		// 1060 = duplicate column
		if !(errors.As(err, &mysqlErr) && mysqlErr.Number == 1060) {
			return fmt.Errorf("failed to add payment_plans.status: %w", err)
		}
	}
	return nil
}

func ensurePaymentPlanIndexes(db *sqlx.DB) error {
	query := 
	`
	CREATE INDEX idx_payment_plan_id_company_updated_at
	ON payment_plans(id_company, updated_at)
	`

	_, err := db.Exec(query)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		// 1061 = duplicate key
		if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1061) {
			return nil
		}
		return fmt.Errorf("failed to create payment_plans indexes: %w", err)
	}

    return nil
}


func ensureInstallments(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS installments(
		id_installment INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		id_payment_plan INT NOT NULL,

		installment_number INT NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		due_date DATE NOT NULL,
		paid_at DATE DEFAULT NULL,

		observations VARCHAR(2000),

		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

		CONSTRAINT unique_id_payment_plan_installment_number UNIQUE (id_payment_plan, installment_number),

		FOREIGN KEY (id_payment_plan) REFERENCES payment_plans(id_payment_plan) ON DELETE CASCADE
	)`)
	if err!=nil{
		return fmt.Errorf("failed to create installments table: %w", err)
	}
	return nil
}
// (id_payment_plan, installment_number) ya tiene índice por constraint UNIQUE (unique_id_payment_plan_installment_number)



func ensureIdempotencyKeys(db *sqlx.DB) error{
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS idempotency_keys (	
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,

		idempotency_key VARCHAR(255) NOT NULL,
		request_hash CHAR(64) NOT NULL,

		resource_type VARCHAR(50),
		resource_id INT,

		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NOT NULL,

		UNIQUE KEY unique_idempotency_key (idempotency_key)
	)`)

	if err!=nil{
		return fmt.Errorf("failed to create idempotency keys table: %w", err)
	}
	return nil

}

func ensureIdempotencyKeysIndexes(db *sqlx.DB) error{
	// idempotency_key ya tiene índice por constraint UNIQUE unique_idempotency_key
	query := `
	CREATE INDEX idx_expires_at
	ON idempotency_keys(expires_at)
	`

	_, err := db.Exec(query)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		// 1061 = duplicate key
		if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1061) {
			return nil
		}
		return fmt.Errorf("failed to create idempotency keys indexes: %w", err)
	}


    return nil
}


func ensureUsers(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users(
		id_user INT PRIMARY KEY AUTO_INCREMENT NOT NULL,

		username VARCHAR(20) NOT NULL,
		password_hash VARCHAR(100) NOT NULL,

		admin BOOLEAN NOT NULL,
		resource_roles JSON,

		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

		CONSTRAINT unique_username UNIQUE (username)
	)`)
	if err!=nil{
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

func ensureRefreshTokens(db *sqlx.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS refresh_tokens(
		id INT PRIMARY KEY AUTO_INCREMENT NOT NULL,
		id_user INT NOT NULL,
		token_hash CHAR(64) NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		revoked_at TIMESTAMP DEFAULT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT unique_token_hash UNIQUE (token_hash),
		FOREIGN KEY (id_user) REFERENCES users(id_user) ON DELETE CASCADE
	)`)
	if err != nil {
		return fmt.Errorf("failed to create refresh_tokens table: %w", err)
	}
	return nil
}

func ensureRefreshTokenIndexes(db *sqlx.DB) error {
	queries := []string{
		`
	CREATE INDEX idx_refresh_token_id_user
	ON refresh_tokens(id_user)
	`,
		`
	CREATE INDEX idx_refresh_token_expires_at
	ON refresh_tokens(expires_at)
	`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1061 {
				continue
			}
			return fmt.Errorf("failed to create refresh_tokens indexes: %w", err)
		}
	}

	return nil
}

func hasUsers(db *sqlx.DB) (bool, error) {
    var exists bool
    err := db.QueryRow(
        "SELECT EXISTS (SELECT 1 FROM users)",
    ).Scan(&exists)

    return exists, err
}

func ensureDefaultAdmin(db *sqlx.DB) error {
	byteHash, err := bcrypt.GenerateFromPassword([]byte("demo1234"), 14)
	if err != nil {
		return apperrors.NewInternalError(fmt.Errorf("failed to generate hash table: %w", err), "")	
	}
	strHash := string(byteHash)

	resourceRoles := map[string]any{
		"company": "editor",
		"member":  "editor",
	}

	resourceRolesJSON, err := json.Marshal(resourceRoles)
	if err!=nil{
		return apperrors.NewInternalError(fmt.Errorf("failed to marshal resource roles: %w", err), "")
	}

	_, err = db.Exec("INSERT INTO users (username, password_hash, admin, resource_roles) VALUES ('admin', ?, true, ?)", strHash, resourceRolesJSON)
	if err!=nil{
		return fmt.Errorf("failed to create default admin: %w", err)
	}
	return nil
}






