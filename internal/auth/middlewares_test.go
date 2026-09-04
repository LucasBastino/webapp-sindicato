package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authadapters "github.com/LucasBastino/webapp-sindicato/internal/auth/adapters"
	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/gofiber/fiber/v2"
)

func TestVerifyToken_SetsSessionRefreshTokenFromCookie(t *testing.T) {
	gen := authadapters.NewJwtTokenGenerator(authadapters.NewJwtNormalizer())
	secret := "0123456789abcdef0123456789abcdef"
	svc := &AuthService{
		tokenGenerator: gen,
		cfg: config.AuthConfig{
			JWTSecret:      secret,
			AccessTokenTTL: time.Hour,
		},
	}
	mw := &AuthMiddleware{authService: svc}

	claims := authdomain.AuthClaims{
		Sub:           1,
		Username:      "user",
		Admin:         false,
		ResourceRoles: map[string]any{"member": "viewer"},
		Exp:           time.Now().Add(time.Hour),
		Iat:           time.Now(),
	}
	accessToken, err := gen.Create(secret, claims)
	if err != nil {
		t.Fatalf("Create access token: %v", err)
	}

	var gotRefresh string
	app := fiber.New()
	app.Use(mw.VerifyToken)
	app.Get("/", func(c *fiber.Ctx) error {
		if tok, ok := c.Locals(SessionRefreshTokenLocal).(string); ok {
			gotRefresh = tok
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: accessToken})
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "request-refresh-token"})

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if gotRefresh != "request-refresh-token" {
		t.Fatalf("expected session refresh from cookie, got %q", gotRefresh)
	}
}
