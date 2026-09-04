package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders sets baseline browser hardening headers.
// HSTS is intentionally omitted (defer until production HTTPS).
func SecurityHeaders() fiber.Handler {
	const csp = "default-src 'self'; " +
		"base-uri 'self'; " +
		"form-action 'self'; " +
		"frame-ancestors 'none'; " +
		"object-src 'none'; " +
		"script-src 'self' 'unsafe-inline'; " +
		"style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; " +
		"font-src 'self' https://fonts.gstatic.com; " +
		"img-src 'self' data:; " +
		"connect-src 'self'"

	return func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", csp)
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Set("Cross-Origin-Opener-Policy", "same-origin")
		return c.Next()
	}
}
