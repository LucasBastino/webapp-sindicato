package backup

import (
	"fmt"

	"github.com/JamesStewy/go-mysqldump"
	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/jmoiron/sqlx"
)

type BackUpRepository struct{
	db *sqlx.DB
}

func NewBackUpRepository(db *sqlx.DB) *BackUpRepository{
	return &BackUpRepository{db: db}
}

func (r BackUpRepository) InitDumper(dumpDir, dumpFilenameFormat string) (*mysqldump.Dumper, error) {
	// mysqldump needs the *sql.DB, and *sqlx.DB.DB returns *sql.DB
	dumper, err := mysqldump.Register(r.db.DB, dumpDir, dumpFilenameFormat)
	if err != nil {
		return nil, apperrors.NewDatabaseError(fmt.Errorf("failed to register dump configuration: %w", err), "")	
	}
	return dumper, nil
}
																											