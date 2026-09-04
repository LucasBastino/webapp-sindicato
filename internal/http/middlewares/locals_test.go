package middlewares

import (
	"io"
	"net/http/httptest"
	"testing"

	userauthinfo "github.com/LucasBastino/webapp-sindicato/internal/features/user/authinfo"
	"github.com/gofiber/fiber/v2"
)

func TestUserIDFromLocals(t *testing.T) {
	app := fiber.New()

	t.Run("absent", func(t *testing.T) {
		app.Get("/absent", func(c *fiber.Ctx) error {
			if got := userIDFromLocals(c); got != "" {
				t.Fatalf("got %q, want empty", got)
			}
			return c.SendStatus(fiber.StatusOK)
		})
		req := httptest.NewRequest("GET", "/absent", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test: %v", err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	})

	t.Run("present", func(t *testing.T) {
		app.Get("/present", func(c *fiber.Ctx) error {
			c.Locals("userAuthInfo", userauthinfo.UserAuthInfo{UserID: 42})
			if got := userIDFromLocals(c); got != "42" {
				t.Fatalf("got %q, want %q", got, "42")
			}
			return c.SendStatus(fiber.StatusOK)
		})
		req := httptest.NewRequest("GET", "/present", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test: %v", err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	})

	t.Run("zero", func(t *testing.T) {
		app.Get("/zero", func(c *fiber.Ctx) error {
			c.Locals("userAuthInfo", userauthinfo.UserAuthInfo{UserID: 0})
			if got := userIDFromLocals(c); got != "" {
				t.Fatalf("got %q, want empty", got)
			}
			return c.SendStatus(fiber.StatusOK)
		})
		req := httptest.NewRequest("GET", "/zero", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test: %v", err)
		}
		defer resp.Body.Close()
		_, _ = io.Copy(io.Discard, resp.Body)
	})
}
