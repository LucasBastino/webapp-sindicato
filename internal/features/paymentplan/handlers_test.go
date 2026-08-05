package paymentplan

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRedirectToPaymentPlanHTMX(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		return redirectToPaymentPlan(c, 42, fiber.StatusCreated)
	})

	req := httptest.NewRequest(fiber.MethodPost, "/", nil)
	req.Header.Set("HX-Request", "true")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("expected status %d, got %d", fiber.StatusCreated, res.StatusCode)
	}
	if location := res.Header.Get("HX-Redirect"); location != "/payment_plans/42" {
		t.Fatalf("expected HX-Redirect to plan, got %q", location)
	}
}

func TestRedirectToPaymentPlanHTTP(t *testing.T) {
	app := fiber.New()
	app.Post("/", func(c *fiber.Ctx) error {
		return redirectToPaymentPlan(c, 42, fiber.StatusCreated)
	})

	req := httptest.NewRequest(fiber.MethodPost, "/", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if res.StatusCode != fiber.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", fiber.StatusSeeOther, res.StatusCode)
	}
	if location := res.Header.Get("Location"); location != "/payment_plans/42" {
		t.Fatalf("expected redirect to plan, got %q", location)
	}
}
