package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/gofiber/fiber/v2"
)

func createAccessCookie(token string, ttl time.Duration, secure bool) fiber.Cookie {
	return fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
		Expires:  time.Now().Add(ttl),
	}
}

func createRefreshCookie(token string, ttl time.Duration, secure bool) fiber.Cookie {
	return fiber.Cookie{
		Name:     "refresh_token",
		Value:    token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
		Expires:  time.Now().Add(ttl),
	}
}

func clearAccessCookie(secure bool) fiber.Cookie {
	return fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	}
}

func clearRefreshCookie(secure bool) fiber.Cookie {
	return fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Lax",
	}
}

func clearCookies(c *fiber.Ctx, secure bool) {
	refreshCookie := clearRefreshCookie(secure)
	c.Cookie(&refreshCookie)

	accessCookie := clearAccessCookie(secure)
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
