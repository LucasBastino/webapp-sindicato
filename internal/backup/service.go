package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
)

type BackUpService struct {
	repo *BackUpRepository
}

func NewBackUpService(repo *BackUpRepository) *BackUpService {
	return &BackUpService{repo: repo}
}

func (s BackUpService) BackUp() error {
	// esto pasarlo por un .env, tanto la direccion como el nombre de la db
	dir := backupDir()

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return apperrors.NewInternalError(fmt.Errorf("failed to create backup directory: %w", err), "")
	}

	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", "sindicatoDB")

	// Register database with mysqldump.
	dumper, err := s.repo.InitDumper(dir, dumpFilenameFormat)
	if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}
	
	// Dump database to file.
	_, err = dumper.Dump()
	if err != nil {
		return apperrors.NewDatabaseError(fmt.Errorf("failed to dump: %w", err), "")
	}

	// // Close dumper, connected database and file stream.
	err = dumper.Close()
	if err!=nil{
		return apperrors.NewDatabaseError(fmt.Errorf("failed to close dumper: %w", err), "")
	}
	return err
}

func (s BackUpService) CleanupOldBackups() error {
	const retentionDays = 60

	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	dir := backupDir()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return apperrors.NewInternalError(fmt.Errorf("failed to read backups directory: %w", err), "")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// backup-2026-06-16.sql
		if !strings.HasPrefix(name, "backup-") || !strings.HasSuffix(name, ".sql") {
			continue
		}

		dateStr := strings.TrimSuffix(strings.TrimPrefix(name, "backup-"), ".sql")

		backupDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// ignoramos nombres inesperados
			continue
		}

		if backupDate.Before(cutoff) {
			path := filepath.Join(dir, name)

			err := os.Remove(path)
			if err != nil {
				return apperrors.NewInternalError(fmt.Errorf("failed to remove old backup %q: %w", name, err), "")
			}
		}
	}

	return nil
}