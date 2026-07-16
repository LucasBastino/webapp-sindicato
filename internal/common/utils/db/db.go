package dbUtils

import (
	"database/sql"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
)

func Exists(db *sql.DB, id int, field, table string) (bool, error) {
	row := db.QueryRow("SELECT * FROM ? WHERE ? = ?", table, field, id)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return true, apperrors.NewDatabaseError(fmt.Errorf("failed to check if model exists in %s by dbUtils: %w", table, err), "")
	}
	return exists, nil
}
