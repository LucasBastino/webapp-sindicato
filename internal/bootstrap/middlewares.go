package bootstrap

import (
	"time"

	"github.com/LucasBastino/app-sindicato/internal/auth"
	"github.com/LucasBastino/app-sindicato/internal/http/middlewares"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
	httplogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func registerMiddlewares(app *fiber.App, authMiddleware *auth.AuthMiddleware, logger logger.Logger) {
	// reciben internamente el ctx.Fiber creado en la request y se ejecutan antes de cada endpoint
	app.Use(middlewares.RecoverMiddleware(logger))
	app.Use(httplogger.New())
	app.Use(requestid.New())
	app.Use(middlewares.TimeoutMiddleWare(10 * time.Second))
	app.Use(middlewares.LoggerMiddleware(logger))
	app.Use(authMiddleware.VerifyAtomicLicense)

}
