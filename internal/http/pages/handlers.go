package pages

import (
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type PagesHandler struct{}

func NewPagesHandler() *PagesHandler {
	return &PagesHandler{}
}

func (h *PagesHandler) RenderDashboard(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	pageContext := page.PageContext{
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "dashboard",
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("dashboard-content", fiber.Map{"PageContext": pageContext})
	}
	return c.Render("pages/dashboard", fiber.Map{"PageContext": pageContext})
}

func (h *PagesHandler) RenderSupport(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	pageContext := page.PageContext{
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "reports",
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("support-content", fiber.Map{"PageContext": pageContext})
	}
	return c.Render("pages/support", fiber.Map{"PageContext": pageContext})
}
