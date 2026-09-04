package config

import "time"

type Config struct {
	Database DatabaseConfig
	Auth     AuthConfig
	License  LicenseConfig
	Server	 ServerConfig
	Backup	 BackupConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	// TLS is the normalized go-sql-driver tls DSN value: "" (disabled), "true", or "skip-verify".
	TLS string
}

type AuthConfig struct {
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	// CookieSecure should be true behind HTTPS (prod/demo). Set COOKIE_SECURE=false for local HTTP.
	CookieSecure bool
}

type LicenseConfig struct{
	Path string
	PublicKeyPath string
}

type ServerConfig struct{
	Port string
	Host string
}

type BackupConfig struct{
	Dir string
}