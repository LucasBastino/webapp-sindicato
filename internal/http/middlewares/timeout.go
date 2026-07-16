package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TimeoutMiddleWare(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), timeout)
		defer cancel() // para liberar recursos cuando el request termina

		c.SetUserContext(ctx)
		return c.Next()
	}
}
