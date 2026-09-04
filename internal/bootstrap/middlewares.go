package bootstrap

import (
	"strings"
	"time"

	"github.com/LucasBastino/webapp-sindicato/internal/auth"
	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/LucasBastino/webapp-sindicato/internal/http/middlewares"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	httplogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func registerMiddlewares(app *fiber.App, authMiddleware *auth.AuthMiddleware, logger logger.Logger, authCfg config.AuthConfig) {
	app.Use(middlewares.RecoverMiddleware(logger))
	app.Use(middlewares.SecurityHeaders())
	app.Use(httplogger.New())
	app.Use(requestid.New())
	app.Use(middlewares.TimeoutMiddleWare(10 * time.Second))
	app.Use(middlewares.LoggerMiddleware(logger))
	app.Use(authMiddleware.VerifyAtomicLicense)
	app.Use(csrf.New(csrf.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.HasPrefix(c.Path(), "/static")
		},
		CookieName:     "csrf_",
		CookiePath:     "/",
		CookieSecure:   authCfg.CookieSecure,
		CookieHTTPOnly: false, // JS reads cookie for HTMX/forms (double-submit)
		CookieSameSite: "Lax",
		Expiration:     12 * time.Hour,
		ContextKey:     "csrf",
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
}
