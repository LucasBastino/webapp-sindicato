package httpUtils

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/common/page"
	userauthinfo "github.com/LucasBastino/webapp-sindicato/internal/features/user/authinfo"
	"github.com/gofiber/fiber/v2"
)


func GetUserAuthInfo(c *fiber.Ctx) (userauthinfo.UserAuthInfo, error) {
	claims, ok := c.Locals("userAuthInfo").(userauthinfo.UserAuthInfo)
	if !ok{
		return userauthinfo.UserAuthInfo{}, apperrors.NewUnauthorizedError(errors.New("invalid auth claims type"), "")
	}
	// return userauthinfo.UserAuthInfo{
	// 	Admin: claims.Admin,
	// 	ResourceRoles: claims.ResourceRoles,
	// }, nil
	return claims, nil
}

func GetIDByParam(c *fiber.Ctx, param string) (int, error) {
	idStr := c.Params(param)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, apperrors.NewBadRequestError(fmt.Errorf("failed to parse id: %w", err), "")
	}

	return id, nil
}

func GetPageByQueryParam(c *fiber.Ctx) int {
	page := c.QueryInt("page")
	if page <= 1 {
		return 1
	} else {
		return page
	}
}

func GetStatusFilters(c *fiber.Ctx) page.StatusFilters {
	active := c.Query("show-active")
	if active == "" {
		active = c.FormValue("show-active")
	}

	inactive := c.Query("show-inactive")
	if inactive == "" {
		inactive = c.FormValue("show-inactive")
	}

	deleted := c.Query("show-deleted")
	if deleted == "" {
		deleted = c.FormValue("show-deleted")
	}

	if active == "" && inactive == "" && deleted == "" {
		return page.StatusFilters{ShowActive: true}
	}

	return page.StatusFilters{
		ShowActive:   active == "true",
		ShowInactive: inactive == "true",
		ShowDeleted:  deleted == "true",
	}
}

func GetSearchKey(c *fiber.Ctx) string {
	key := c.Query("search-key")
	if key == "" {
		key = c.FormValue("search-key")
	}
	return key
}

// SessionRefreshToken returns the active refresh token for the current request.
// After middleware rotation it may differ from the refresh_token request cookie.
func SessionRefreshToken(c *fiber.Ctx) string {
	if t, ok := c.Locals("sessionRefreshToken").(string); ok && t != "" {
		return t
	}
	return c.Cookies("refresh_token")
}
