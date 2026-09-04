package auth

import (
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/gofiber/fiber/v2"
)

func writeSessionCookies(c *fiber.Ctx, cfg config.AuthConfig, accessToken, refreshToken string) {
	accessCookie := createAccessCookie(accessToken, cfg.AccessTokenTTL, cfg.CookieSecure)
	c.Cookie(&accessCookie)
	refreshCookie := createRefreshCookie(refreshToken, cfg.RefreshTokenTTL, cfg.CookieSecure)
	c.Cookie(&refreshCookie)
}

// HasValidSession reports whether the request carries a valid access or refresh session.
// When refresh succeeds, new cookies are written on c.
func (s *AuthService) HasValidSession(c *fiber.Ctx) bool {
	ctx := c.UserContext()

	token := c.Cookies("access_token")
	claims, err := s.VerifyAccessToken(token)
	if err == nil && claims != nil && !claims.Exp.Before(time.Now()) {
		return true
	}

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return false
	}

	newAccessToken, newRefreshToken, _, err := s.RefreshAccessToken(ctx, refreshToken)
	if err != nil {
		return false
	}

	writeSessionCookies(c, s.cfg, newAccessToken, newRefreshToken)
	return true
}
