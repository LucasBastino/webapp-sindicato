package middlewares

import (
	"fmt"
	"runtime/debug"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/http/errorhandler"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
)

func RecoverMiddleware(logger logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			p := recover()
			if p != nil {
				requestID, ok := c.Locals("requestid").(string)
				if !ok {
					requestID = ""
				}

				userID := userIDFromLocals(c)
				logger.Error("request panicked",
					"request_id", requestID,
					"user_id", userID,
					"method", c.Method(),
					"path", c.Path(),
					"status_code", 500,
					"panic_value", p,
					"stack", debug.Stack(),
				)
				err := apperrors.NewInternalError(fmt.Errorf("panic: %v", p), "")
				_ = errorhandler.HandleError(c,  err)
				return
			}
		}()
		// si no hay panic, pasa directamente al proximo middleware o handler
		return c.Next()
	}
}
