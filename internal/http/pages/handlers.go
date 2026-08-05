package pages

import (
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type PagesHandler struct {
	dashboard *DashboardService
}

func NewPagesHandler(dashboard *DashboardService) *PagesHandler {
	return &PagesHandler{dashboard: dashboard}
}

type dashboardPageData struct {
	Summary     DashboardSummary
	PageContext page.PageContext
}

func (h *PagesHandler) RenderDashboard(c *fiber.Ctx) error {
	ctx := c.UserContext()
	
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	summary, err := h.dashboard.GetSummary(ctx)
	if err != nil {
		return err
	}

	pageData := dashboardPageData{
		Summary: summary,
		PageContext: page.PageContext{
			UserAuthInfo:  userAuthInfo,
			ActiveSection: "dashboard",
		},
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("dashboard-content", pageData)
	}
	return c.Render("pages/dashboard", pageData)
}

func (h *PagesHandler) RenderReports(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	pageData := fiber.Map{
		"PageContext": page.PageContext{
			UserAuthInfo:  userAuthInfo,
			ActiveSection: "reports",
		},
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("reports-content", pageData)
	}
	return c.Render("pages/reports", pageData)
}

func (h *PagesHandler) RenderSupport(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	pageContext := page.PageContext{
		UserAuthInfo: userAuthInfo,
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("support-modal", fiber.Map{"PageContext": pageContext})
	}
	return c.Render("pages/support", fiber.Map{"PageContext": pageContext})
}
