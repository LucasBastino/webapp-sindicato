package payment

import (
	"fmt"
	"strconv"
	"time"

	authports "github.com/LucasBastino/webapp-sindicato/internal/auth/ports"
	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/webapp-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	service    *PaymentService
	normalizer authports.ClaimsNormalizer
}

func NewPaymentHandler(service *PaymentService, normalizer authports.ClaimsNormalizer) *PaymentHandler {
	return &PaymentHandler{
		service:    service,
		normalizer: normalizer,
	}
}

func (h *PaymentHandler) RenderModal(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	payment, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}

	res := toResponse(*payment)
	pageContext := page.PageContext{
		Mode:          "edit",
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "companies",
	}
	pageData := pageData{
		Payment:     res,
		CompanyID:   payment.CompanyID,
		PageContext: pageContext,
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("payment-modal", pageData)
	}
	return c.Render("payment/payment", pageData)
}

func (h *PaymentHandler) RenderGrid(c *fiber.Ctx) error {
	ctx := c.UserContext()
	year := c.Query("year")

	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err != nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	companyName, number, address, phone, err := h.service.GetCompanyNavDetails(ctx, companyID)
	if err != nil {
		return err
	}

	totalRows, err := h.service.Count(ctx, companyID)
	if err != nil {
		return err
	}

	pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "edit", ActiveSection: "companies"}
	data := gridPageData{
		CompanyID:   companyID,
		CompanyName: companyName,
		CompanyNav: page.NewCompanyNav(companyID, companyName, "payments", userAuthInfo.CanView("member")).
			WithDetails(number, address, phone),
		PageContext: pageContext,
	}

	renderPayments := func(d gridPageData) error {
		if c.Get("HX-Request") == "true" {
			return c.Render("payments-content", d)
		}
		return c.Render("payment/payments", d)
	}

	if totalRows == 0 {
		return renderPayments(data)
	}

	years, err := h.service.ListPaymentYears(ctx, companyID)
	if err != nil {
		return err
	}
	if len(years) == 0 {
		return renderPayments(data)
	}

	data.Years = years
	yearInt, err := resolvePaymentYear(year, years, time.Now())
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse year: %w", err), "")
	}

	payments, err := h.service.List(ctx, companyID, yearInt)
	if err != nil {
		return err
	}
	data.Payments = toGridResponses(payments)
	data.Stats = buildGridStats(payments)
	data.Year = yearInt
	return renderPayments(data)
}

func (h *PaymentHandler) RenderOverdue(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	groups, err := h.service.ListOverdue(ctx)
	if err != nil {
		return err
	}

	pageContext := page.PageContext{
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "reports",
	}
	data := overduePageData{
		Groups:      groups,
		PageContext: pageContext,
	}
	if len(groups) == 0 {
		data.EmptyState = page.EmptyState{
			Icon:        "alert-triangle",
			Title:       "No hay aportes vencidos",
			Description: "No hay aportes vencidos para mostrar.",
		}
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("overdue-payments-content", data)
	}
	return c.Render("payment/overdue_payments", data)
}

func (h *PaymentHandler) CreatePayments(c *fiber.Ctx) error {
	ctx := c.UserContext()
	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err != nil {
		return err
	}
	return h.service.CreateRemainingPaymentsByID(ctx, nil, companyID)
}

func (h *PaymentHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	var req request
	err = c.BodyParser(&req)
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		payment, err := h.service.Get(ctx, id)
		if err != nil {
			return err
		}
		res, err := mergetoResponse(*payment, req)
		if err != nil {
			return err
		}
		pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "companies"}
		pageData := pageData{Payment: res, CompanyID: payment.CompanyID, PageContext: pageContext, Errors: errorMap}
		c.Set("HX-Retarget", "#app-modal-container")
		c.Set("HX-Reswap", "innerHTML")
		status := fiber.StatusBadRequest
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		return c.Status(status).Render("payment-modal", pageData)
	}

	payment, err := toModel(req)
	if err != nil {
		return err
	}
	err = h.service.Update(ctx, id, payment)
	if err != nil {
		return err
	}

	// Close modal first; refresh the underlying payments/overdue view after settle.
	c.Set("HX-Trigger-After-Settle", "refreshPayments")
	return c.SendString("")
}

func resolvePaymentYear(requestedYear string, availableYears []int, now time.Time) (int, error) {
	if requestedYear != "" && requestedYear != "0" {
		parsed, err := strconv.Atoi(requestedYear)
		if err != nil {
			return 0, err
		}
		return parsed, nil
	}

	currentYear := now.Year()
	for _, y := range availableYears {
		if y == currentYear {
			return currentYear, nil
		}
	}

	if len(availableYears) > 0 {
		return availableYears[0], nil
	}

	return currentYear, nil
}
