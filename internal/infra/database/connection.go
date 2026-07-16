package database

import (
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// connection to mysql database and return the *sqlx.DB
func Connect(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, apperrors.NewDatabaseError(fmt.Errorf("failed to connect to mysql database: %w", err), "")
	}

	// err = db.Ping()
	// if err != nil {
	// 	return nil, apperrors.NewDatabaseError(fmt.Errorf("failed to ping database: %w", err), "")
	// }
	return db, nil
}