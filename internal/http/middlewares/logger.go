package middlewares

import (
	"errors"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
)

func LoggerMiddleware(logger logger.Logger) fiber.Handler {
	// no loggea el GET a favicon
	return func(c *fiber.Ctx) error {
		if c.Path() == "/favicon.ico" {
        return c.Next()
    }

		// en fiber v2 el middleware oficial de requestid lo guarda automaticamente en locals como "requestid"
		requestID, ok := c.Locals("requestid").(string)
		if !ok {
			requestID = ""
		}

		userID, ok := c.Locals("userID").(string)
		if !ok {
			userID = ""
		}

		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		// los inicializo antes para mas seguridad
		errorType := "internal"
		statusCode := fiber.StatusInternalServerError

		if err != nil {
			var appErr *apperrors.AppError
			// errors.As comprueba si hay un error del tipo CustomError dentro de err,
			// y si lo hay, lo asigna a la variable customErr
			if errors.As(err, &appErr) {
				errorType = appErr.Type
				statusCode = appErr.StatusCode
			}

			logger.Error("request_failed",
				"request_id", requestID,
				"user_id", userID,
				"method", c.Method(),
				"path", c.Path(),
				"status_code", statusCode,
				"duration_ms", duration.Milliseconds(),
				"error_type", errorType,
				"internal_msg", err.Error(),
			)

			return err
		}

		logger.Info("request_completed",
			"request_id", requestID,
			"user_id", userID,
			"method", c.Method(),
			"path", c.Path(),
			// este es un getter ↓, c.Status es un setter
			"status_code", c.Response().StatusCode(),
			"duration_ms", duration.Milliseconds(),
		)

		return nil
	}
}