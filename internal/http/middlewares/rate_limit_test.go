package middlewares_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LucasBastino/webapp-sindicato/internal/http/middlewares"
	"github.com/gofiber/fiber/v2"
)

func TestLoginRateLimiter_Returns429AfterMax(t *testing.T) {
	app := fiber.New()
	app.Post("/login", middlewares.LoginRateLimiter(), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("request %d: status=%d body=%s", i+1, resp.StatusCode, body)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/login", nil)
	req.RemoteAddr = "1.2.3.4:1234"
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 429, got %d body=%s", resp.StatusCode, body)
	}
}

func TestVerifyLicenseRateLimiter_Returns429AfterMax(t *testing.T) {
	app := fiber.New()
	app.Post("/verifyLicense", middlewares.VerifyLicenseRateLimiter(), func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodPost, "/verifyLicense", nil)
		req.RemoteAddr = "5.6.7.8:5678"
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("request %d: status=%d body=%s", i+1, resp.StatusCode, body)
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/verifyLicense", nil)
	req.RemoteAddr = "5.6.7.8:5678"
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 429, got %d body=%s", resp.StatusCode, body)
	}
}
