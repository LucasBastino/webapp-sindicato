package config

import (
	"fmt"
	"strings"
)

const minJWTSecretLen = 32

// Validate fails fast when critical secrets/config required to boot securely are missing.
func (c Config) Validate() error {
	var missing []string

	if strings.TrimSpace(c.Auth.JWTSecret) == "" {
		missing = append(missing, "JWT_SECRET")
	} else if len(c.Auth.JWTSecret) < minJWTSecretLen {
		return fmt.Errorf("JWT_SECRET must be at least %d characters", minJWTSecretLen)
	}

	if strings.TrimSpace(c.Database.User) == "" {
		missing = append(missing, "DB_USER")
	}
	if strings.TrimSpace(c.Database.Host) == "" {
		missing = append(missing, "DB_HOST")
	}
	if strings.TrimSpace(c.Database.Port) == "" {
		missing = append(missing, "DB_PORT")
	}
	if strings.TrimSpace(c.Database.Name) == "" {
		missing = append(missing, "DB_NAME")
	}

	if strings.TrimSpace(c.License.Path) == "" {
		missing = append(missing, "LICENSE_PATH")
	}
	if strings.TrimSpace(c.License.PublicKeyPath) == "" {
		missing = append(missing, "PUBLIC_KEY_PATH")
	}

	if strings.TrimSpace(c.Server.Port) == "" {
		missing = append(missing, "APP_PORT")
	}
	if strings.TrimSpace(c.Server.Host) == "" {
		missing = append(missing, "APP_HOST")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required env: %s", strings.Join(missing, ", "))
	}

	switch c.Database.TLS {
	case "", "true", "skip-verify":
		// ok
	default:
		return fmt.Errorf("DB_TLS must be false, true, or skip-verify")
	}

	return nil
}
