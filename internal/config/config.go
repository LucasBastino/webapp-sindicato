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
}

type AuthConfig struct {
	JWTSecret      string
	AccessTokenTTL time.Duration
	RefreshTokenTTL time.Duration
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