package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/gofiber/fiber/v2"
)

func createAccessCookie(token string, ttl time.Duration) fiber.Cookie {
	return fiber.Cookie{
		Name:		"access_token",
		Value:      token,
		Path:       "/",
		HTTPOnly:   true,
		Secure:     false,
		SameSite:   "Lax",
		Expires:	time.Now().Add(ttl),
		// SessionOnly: true,
		// para subir a un dominio
		// Secure:   true,
		// SameSite: "None",
	}
}

func createRefreshCookie(token string, ttl time.Duration) fiber.Cookie {
	return fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,        // true en producción con HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(ttl),
	}
}

func clearAccessCookie() fiber.Cookie {
	return fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expired 1 hour ago
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		// para subir a un dominio
		// Secure:   true,
		// SameSite: "None",
	}
}

func clearRefreshCookie() fiber.Cookie {
	return fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	}
}

func clearCookies(c *fiber.Ctx) {
	refreshCookie := clearRefreshCookie()
	c.Cookie(&refreshCookie)

	accessCookie := clearAccessCookie()
	c.Cookie(&accessCookie)
}

// Para el refresh token se usa un hash determinístico, no como el bcrypt que es costoso, ese se usa para constraseñas 
func hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}

func generateRandomToken() (string, error) {
    bytes := make([]byte, 32)
    _, err := rand.Read(bytes)
    if err != nil {
        return "", apperrors.NewInternalError(fmt.Errorf("failed to generate refresh token: %w", err), "")
    }
    return hex.EncodeToString(bytes), nil
}