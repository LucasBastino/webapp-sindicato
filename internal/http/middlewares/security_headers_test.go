package middlewares

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(SecurityHeaders())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	want := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Permissions-Policy":        "camera=(), microphone=(), geolocation=()",
		"Cross-Origin-Opener-Policy": "same-origin",
	}
	for header, value := range want {
		if got := resp.Header.Get(header); got != value {
			t.Errorf("%s: got %q, want %q", header, got, value)
		}
	}

	csp := resp.Header.Get("Content-Security-Policy")
	if csp == "" {
		t.Fatal("Content-Security-Policy missing")
	}
	for _, part := range []string{
		"default-src 'self'",
		"frame-ancestors 'none'",
		"script-src 'self' 'unsafe-inline'",
		"https://fonts.googleapis.com",
		"https://fonts.gstatic.com",
	} {
		if !strings.Contains(csp, part) {
			t.Errorf("CSP missing %q; got %q", part, csp)
		}
	}
	if strings.Contains(csp, "unpkg.com") {
		t.Errorf("CSP must not allow unpkg; got %q", csp)
	}

	if resp.Header.Get("Strict-Transport-Security") != "" {
		t.Fatal("HSTS must not be set")
	}
}
