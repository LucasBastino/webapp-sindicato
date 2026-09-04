package config

import (
	"os"
	"strings"
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
			TLS:      parseDBTLS(os.Getenv("DB_TLS")),
		},
		Auth: AuthConfig{
			JWTSecret:       os.Getenv("JWT_SECRET"),
			AccessTokenTTL:  5 * time.Minute,
			RefreshTokenTTL: 7 * 24 * time.Hour,
			CookieSecure:    parseCookieSecure(os.Getenv("COOKIE_SECURE")),
		},
		License: LicenseConfig{
			Path:          os.Getenv("LICENSE_PATH"),
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

// parseCookieSecure defaults to true (assume HTTPS). Explicit false/0/no disables for local HTTP.
func parseCookieSecure(raw string) bool {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "false", "0", "no":
		return false
	default:
		return true
	}
}

// parseDBTLS returns the normalized tls DSN value.
// Unknown values are returned lowercased so Validate can fail fast.
func parseDBTLS(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", "false", "0", "no":
		return ""
	case "true", "1", "yes":
		return "true"
	case "skip-verify":
		return "skip-verify"
	default:
		return strings.ToLower(strings.TrimSpace(raw))
	}
}
