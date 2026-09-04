package bootstrap_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
)

func TestCSRF_RejectsMutationWithoutToken(t *testing.T) {
	app := fiber.New()
	app.Use(csrf.New(csrf.Config{
		CookieName:     "csrf_",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: false,
		CookieSameSite: "Lax",
		Extractor: func(c *fiber.Ctx) (string, error) {
			if token := c.Get(csrf.HeaderName); token != "" {
				return token, nil
			}
			if token := c.FormValue("csrf"); token != "" {
				return token, nil
			}
			return "", csrf.ErrTokenNotFound
		},
	}))
	app.Get("/ping", func(c *fiber.Ctx) error { return c.SendString("ok") })
	app.Post("/mutate", func(c *fiber.Ctx) error { return c.SendString("mutated") })

	getReq := httptest.NewRequest(http.MethodGet, "/ping", nil)
	getResp, err := app.Test(getReq)
	if err != nil {
		t.Fatal(err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET status=%d", getResp.StatusCode)
	}

	cookie := csrfCookieFromResponse(getResp)
	if cookie == "" {
		t.Fatal("expected csrf_ cookie on GET")
	}

	badReq := httptest.NewRequest(http.MethodPost, "/mutate", nil)
	badReq.AddCookie(&http.Cookie{Name: "csrf_", Value: cookie})
	badResp, err := app.Test(badReq)
	if err != nil {
		t.Fatal(err)
	}
	if badResp.StatusCode != http.StatusForbidden {
		t.Fatalf("POST without token status=%d want 403", badResp.StatusCode)
	}

	form := strings.NewReader("csrf=" + cookie)
	okReq := httptest.NewRequest(http.MethodPost, "/mutate", form)
	okReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	okReq.AddCookie(&http.Cookie{Name: "csrf_", Value: cookie})
	okResp, err := app.Test(okReq)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(okResp.Body)
	if okResp.StatusCode != http.StatusOK {
		t.Fatalf("POST with token status=%d body=%s", okResp.StatusCode, body)
	}
}

func csrfCookieFromResponse(resp *http.Response) string {
	for _, c := range resp.Cookies() {
		if c.Name == "csrf_" {
			return c.Value
		}
	}
	for _, h := range resp.Header.Values("Set-Cookie") {
		if strings.HasPrefix(h, "csrf_=") {
			part := strings.SplitN(h, ";", 2)[0]
			return strings.TrimPrefix(part, "csrf_=")
		}
	}
	return ""
}
