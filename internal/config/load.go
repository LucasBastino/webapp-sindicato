package config

import (
	"os"
	"time"
)

func Load() Config {
	return Config{
		Database: DatabaseConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
		},
		Auth: AuthConfig{
			JWTSecret: os.Getenv("JWT_SECRET"),
			AccessTokenTTL: 5 * time.Minute,
			RefreshTokenTTL: 7*24 * time.Hour,
		},
		License: LicenseConfig{
			Path: os.Getenv("LICENSE_PATH"),
			PublicKeyPath: os.Getenv("PUBLIC_KEY_PATH"),
		},
		Server: ServerConfig{
			Port: os.Getenv("APP_PORT"),
			Host: os.Getenv("APP_HOST"),
		},
		Backup: BackupConfig{
			Dir: os.Getenv("BACKUP_DIR"),
		},
	}
}