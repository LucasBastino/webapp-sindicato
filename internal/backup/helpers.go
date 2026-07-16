package backup

import "os"

func backupDir() string {
	dir := os.Getenv("BACKUP_DIR")
	if dir == "" {
		dir = "./backups"
	}
	return dir
}