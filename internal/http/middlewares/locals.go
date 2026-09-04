package middlewares

import (
	"strconv"

	userauthinfo "github.com/LucasBastino/webapp-sindicato/internal/features/user/authinfo"
	"github.com/gofiber/fiber/v2"
)

// userIDFromLocals returns the authenticated user id as a string, or "" if absent.
func userIDFromLocals(c *fiber.Ctx) string {
	info, ok := c.Locals("userAuthInfo").(userauthinfo.UserAuthInfo)
	if !ok || info.UserID == 0 {
		return ""
	}
	return strconv.Itoa(info.UserID)
}
