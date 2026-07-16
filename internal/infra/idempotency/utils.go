package idempotency

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func generateRequestHash(c *fiber.Ctx) (string, error) {
	body := c.Body()

	data := fmt.Sprintf(
		"%s:%s:%s",
		c.Method(),
		c.Path(),
		string(body),
	)

	hash := sha256.Sum256([]byte(data))

	return hex.EncodeToString(hash[:]), nil
}