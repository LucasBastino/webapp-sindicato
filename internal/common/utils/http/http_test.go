package httpUtils

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func TestSessionRefreshToken_PrefersLocals(t *testing.T) {
	app := fiber.New()
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetCookie("refresh_token", "old-request-token")

	c := app.AcquireCtx(ctx)
	defer app.ReleaseCtx(c)

	c.Locals("sessionRefreshToken", "rotated-token")

	if got := SessionRefreshToken(c); got != "rotated-token" {
		t.Fatalf("expected rotated-token, got %q", got)
	}
}

func TestSessionRefreshToken_FallsBackToCookie(t *testing.T) {
	app := fiber.New()
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetCookie("refresh_token", "cookie-token")

	c := app.AcquireCtx(ctx)
	defer app.ReleaseCtx(c)

	if got := SessionRefreshToken(c); got != "cookie-token" {
		t.Fatalf("expected cookie-token, got %q", got)
	}
}
