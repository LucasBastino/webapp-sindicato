package paymentplan

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/gofiber/fiber/v2"
)

type PaymentPlanHandler struct {
	service *PaymentPlanService
}

func redirectToPaymentPlan(c *fiber.Ctx, id, status int) error {
	location := fmt.Sprintf("/payment_plans/%d", id)
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", location)
		return c.SendStatus(status)
	}
	return c.Redirect(location, fiber.StatusSeeOther)
}

func NewPaymentPlanHandler(service *PaymentPlanService) *PaymentPlanHandler {
	return &PaymentPlanHandler{
		service: service,
	}
}

func (h *PaymentPlanHandler) RenderOverview(c *fiber.Ctx) error {
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	groups, err := h.service.ListAllGrouped(ctx)
	if err != nil {
		return err
	}

	pageContext := page.PageContext{
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "reports",
	}
	data := overviewPageData{
		Groups:      groups,
		PageContext: pageContext,
	}
	if len(groups) == 0 {
		data.EmptyState = page.EmptyState{
			Icon:        "credit-card",
			Title:       "No hay planes de pago",
			Description: "No hay planes de pago para mostrar.",
		}
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("payment-plans-overview-content", data)
	}
	return c.Render("paymentPlan/payment_plans_overview", data)
}

func (h *PaymentPlanHandler) RenderPage(c *fiber.Ctx) error {
	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}
	return h.renderPageByID(c, id)
}

func (h *PaymentPlanHandler) renderPageByID(c *fiber.Ctx, id int) error {
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	detail, err := h.service.GetDetail(ctx, id)
	if err != nil {
		return err
	}

	res := toResponse(*detail)

	pageContext := page.PageContext{
		Mode:          "edit",
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "companies",
	}

	pageData := pageData{
		PaymentPlan: res,
		CompanyID:   detail.CompanyID,
		PageContext: pageContext,
	}
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Push-Url", fmt.Sprintf("/payment_plans/%d", id))
		return c.Render("payment-plan-content", pageData)
	}
	return c.Render("paymentPlan/payment_plan", pageData)
}

func (h *PaymentPlanHandler) RenderTable(c *fiber.Ctx) error {

	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err != nil {
		return err
	}

	return h.renderTableByID(c, companyID)
}

func (h *PaymentPlanHandler) renderTableByID(c *fiber.Ctx, companyID int) error {
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	totalRows, err := h.service.Count(ctx, companyID)
	if err != nil {
		return err
	}

	pageContext := page.PageContext{
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "companies",
	}

	var paymentPlans []tableResponse
	var emptyState page.EmptyState

	if totalRows == 0 {
		if userAuthInfo.CanEdit("company") {
			emptyState = page.NewEmptyState(
				"file-text",
				"planes de pago",
				fmt.Sprintf("/companies/%d/payment_plans/new", companyID),
				"Crear plan",
			)
		} else {
			emptyState = page.NewNoResultsEmptyState("file-text", "planes de pago")
		}
	} else {
		list, err := h.service.List(ctx, companyID)
		if err != nil {
			return err
		}
		paymentPlans = toTableResponses(list)
	}

	tablePageData := tablePageData{
		PaymentPlans: paymentPlans,
		TotalResults: totalRows,
		EmptyState:   emptyState,
		CompanyID:    companyID,
		PageContext:  pageContext,
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("payment-plans-content", tablePageData)
	}
	return c.Render("paymentPlan/payment_plans", tablePageData)
}

func (h *PaymentPlanHandler) RenderAddForm(c *fiber.Ctx) error {
	ctx := c.UserContext()
	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err != nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	overdue, err := h.service.ListOverdueForCompany(ctx, companyID)
	if err != nil {
		return err
	}

	pageContext := page.PageContext{
		Mode:          "add",
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "companies",
	}

	pageData := pageData{
		PaymentPlan:     response{SelectedPaymentIDs: map[int]bool{}},
		OverduePayments: toOverdueOptions(overdue),
		CompanyID:       companyID,
		PageContext:     pageContext,
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("add-payment-plan-content", pageData)
	}
	return c.Render("paymentPlan/add_payment_plan", pageData)
}

func (h *PaymentPlanHandler) Create(c *fiber.Ctx) error {
	record := c.Locals("idempotency_record")
	if record != nil {

		r, ok := record.(*idempotency.Record)
		if !ok {
			return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
		}

		return redirectToPaymentPlan(c, *r.ResourceID, fiber.StatusOK)

	}

	ctx := c.UserContext()
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	var req request
	err = c.BodyParser(&req)
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment plan req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		res := toResponseFromRequest(req)
		companyID, err := strconv.Atoi(req.CompanyID)
		if err != nil {
			return apperrors.NewBadRequestError(fmt.Errorf("invalid company id: %w", err), "")
		}
		overdue, err := h.service.ListOverdueForCompany(ctx, companyID)
		if err != nil {
			return err
		}
		pageContext := page.PageContext{
			UserAuthInfo:  userAuthInfo,
			Mode:          "add",
			ActiveSection: "companies",
		}
		pageData := pageData{
			PaymentPlan:     res,
			OverduePayments: toOverdueOptions(overdue),
			CompanyID:       companyID,
			PageContext:     pageContext,
			Errors:          errorMap,
		}
		status := fiber.StatusBadRequest
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		c.Status(status)
		if c.Get("HX-Request") == "true" {
			return c.Render("add-payment-plan-content", pageData)
		}
		return c.Render("paymentPlan/add_payment_plan", pageData)
	}

	input, err := toCreateInput(req)
	if err != nil {
		return err
	}

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok {
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	id, err := h.service.Create(ctx, input, idempotencyKey)
	if err != nil {
		res := toResponseFromRequest(req)
		companyID, err2 := strconv.Atoi(req.CompanyID)
		if err2 != nil {
			return apperrors.NewBadRequestError(fmt.Errorf("invalid company id: %w", err2), "")
		}
		overdue, listErr := h.service.ListOverdueForCompany(ctx, companyID)
		if listErr != nil {
			return listErr
		}
		pageContext := page.PageContext{
			UserAuthInfo:  userAuthInfo,
			Mode:          "add",
			ActiveSection: "companies",
		}
		pageData := pageData{
			PaymentPlan:     res,
			OverduePayments: toOverdueOptions(overdue),
			CompanyID:       companyID,
			PageContext:     pageContext,
			Errors:          map[string]string{},
		}
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) && appErr.ClientMsg != "" {
			pageData.Errors["form"] = appErr.ClientMsg
		} else {
			pageData.Errors["form"] = "No se pudo crear el plan de pago."
		}
		status := fiber.StatusConflict
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		c.Status(status)
		if c.Get("HX-Request") == "true" {
			return c.Render("add-payment-plan-content", pageData)
		}
		return c.Render("paymentPlan/add_payment_plan", pageData)
	}

	return redirectToPaymentPlan(c, id, fiber.StatusCreated)
}

func (h *PaymentPlanHandler) Update(c *fiber.Ctx) error {
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
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment plan req body: %w", err), "")
	}

	req.trim()
	errorMap := map[string]string{}
	if msg := validateObservationsOnly(req.Observations); msg != "" {
		errorMap["observations"] = msg
	}
	if len(errorMap) > 0 {
		detail, err := h.service.GetDetail(ctx, id)
		if err != nil {
			return err
		}
		res, err := mergetoResponse(*detail, req)
		if err != nil {
			return err
		}
		pageContext := page.PageContext{
			Mode:          "edit",
			UserAuthInfo:  userAuthInfo,
			ActiveSection: "companies",
		}
		pageData := pageData{
			PaymentPlan: res,
			CompanyID:   detail.CompanyID,
			PageContext: pageContext,
			Errors:      errorMap,
		}
		if c.Get("HX-Request") == "true" {
			return c.Status(fiber.StatusOK).Render("payment-plan-content", pageData)
		}
		return c.Status(fiber.StatusBadRequest).Render("paymentPlan/payment_plan", pageData)
	}

	err = h.service.Update(ctx, id, PaymentPlan{Observations: req.Observations})
	if err != nil {
		return err
	}

	return h.RenderPage(c)
}

func (h *PaymentPlanHandler) Cancel(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	err = h.service.Cancel(ctx, id)
	if err != nil {
		return err
	}

	if c.Get("view") == "table" {
		paymentPlan, _, err := h.service.Get(ctx, id)
		if err != nil {
			return err
		}
		return h.renderTableByID(c, paymentPlan.CompanyID)
	}

	return h.renderPageByID(c, id)
}

func (h *PaymentPlanHandler) Restore(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	err = h.service.Restore(ctx, id)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) && appErr.Type == "business" && c.Get("HX-Request") == "true" {
			c.Set("HX-Retarget", "#app-modal-container")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusOK).Render("error/error_modal", fiber.Map{
				"errorMsg": appErr.ClientMsg,
			})
		}
		return err
	}

	return h.renderPageByID(c, id)
}

func (h *PaymentPlanHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	paymentPlan, _, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}

	err = h.service.HardDelete(ctx, id)
	if err != nil {
		return err
	}

	return h.renderTableByID(c, paymentPlan.CompanyID)
}
