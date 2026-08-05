package idempotency

import (
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/gofiber/fiber/v2"
)

type IdempotencyMiddleware struct {
	service *IdempotencyService
}

func NewIdempotencyMiddleware(service *IdempotencyService) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{service: service}
}

const HeaderIdempotencyKey = "Idempotency-Key"

func (m *IdempotencyMiddleware) VerifyIdempotency(c *fiber.Ctx) error {
	key := c.Get(HeaderIdempotencyKey)

	if key == "" {
		return apperrors.NewBadRequestError(errors.New("missing idempotency key"), "")
	}

	hash, err := generateRequestHash(c)
	if err != nil {
		return apperrors.NewInternalError(fmt.Errorf("failed to generate request hash: %w", err), "")
	}

	record, err := m.service.CheckOrCreate(c.Context(),	key, hash)
	if err != nil {
		return err
	}

	if record != nil {
		c.Locals("idempotency_record", record)
	}

	c.Locals("idempotency_key", key)
	return c.Next()
}

