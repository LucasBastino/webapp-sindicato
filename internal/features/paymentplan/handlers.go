package paymentplan

import (
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/gofiber/fiber/v2"
)

type PaymentPlanHandler struct {
	service *PaymentPlanService
	idempotencyService *idempotency.IdempotencyService
}

func NewPaymentPlanHandler(service *PaymentPlanService, idempotencyService *idempotency.IdempotencyService) *PaymentPlanHandler {
	return &PaymentPlanHandler{
		service: service,
		idempotencyService: idempotencyService,
	}
}

func (h *PaymentPlanHandler) RenderPage(c *fiber.Ctx) error {
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil{
		return err
	}
	return h.renderPageByID(c, id)
}

func (h *PaymentPlanHandler) renderPageByID(c *fiber.Ctx, id int) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	paymentPlan, installments, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}

	res := toResponse(*paymentPlan, installments)

	pageContext := page.PageContext{
		Mode: "edit",
		UserAuthInfo: userAuthInfo,
		ActiveSection: "companies",
	}

	pageData := pageData{
		PaymentPlan: res,
		PageContext: pageContext,
	}
	if c.Get("HX-Request") == "true" {
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

func (h *PaymentPlanHandler) renderTableByID(c *fiber.Ctx, companyID int) error{
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	totalRows, err := h.service.Count(ctx, companyID)
	if err != nil {
		return err
	}

	renderPlans := func(data tablePageData) error {
		if c.Get("HX-Request") == "true" {
			return c.Render("payment-plans-content", data)
		}
		return c.Render("paymentPlan/payment_plans", data)
	}

	if totalRows == 0 {
		pageContext := page.PageContext{UserAuthInfo: userAuthInfo, ActiveSection: "companies"}
		tablePageData := tablePageData{PageContext: pageContext}
		return renderPlans(tablePageData)
	}

	paymentPlans, err := h.service.List(ctx, companyID)
	if err != nil {
		return err
	}

	responses := toTableResponses(paymentPlans)

	pageContext := page.PageContext{UserAuthInfo: userAuthInfo, ActiveSection: "companies"}
	tablePageData := tablePageData{
		PaymentPlans: responses,
		PageContext: pageContext,
	}
	return renderPlans(tablePageData)
}

func (h *PaymentPlanHandler) RenderAddForm(c *fiber.Ctx) error{
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	pageContext := page.PageContext{
		Mode: "add",
		UserAuthInfo: userAuthInfo,
		ActiveSection: "companies",
	}

	pageData := pageData{
		PaymentPlan: response{},
		PageContext: pageContext,
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
		if !ok{
			return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
		}

		c.Set("HX-Redirect", fmt.Sprintf("/payment_plans/%d", *r.ResourceID))
		return c.SendStatus(200)

	}
	
	ctx := c.UserContext()
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	var req request
	err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment plan req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		res, err := toResponseFromRequest(req)
		if err!=nil{
			return err
		}		
		PageContext := page.PageContext{
			UserAuthInfo: userAuthInfo,
			Mode:         "add",
			ActiveSection: "companies",
		}
		pageData := pageData{
			PaymentPlan: res,
			PageContext: PageContext,
			Errors:      errorMap,
		}
		return c.Status(fiber.StatusBadRequest).Render("paymentPlan/add_payment_plan", pageData)
	}


	// Si no tiene errores inserto el member en la DB y renderizo el su archivo
	paymentPlan, err := toModel(req)
	if err!=nil{
		return err
	}
	id, err := h.service.Create(ctx, paymentPlan)
	if err!=nil{
		res, err := toResponseFromRequest(req)
		if err!=nil{
			return err
		}
		PageContext := page.PageContext{
			UserAuthInfo: userAuthInfo,
			Mode:         "add",
			ActiveSection: "companies",
		}
		pageData := pageData{
			PaymentPlan: res,
			PageContext: PageContext,
			Errors:      errorMap,
		}
		return c.Status(fiber.StatusConflict).Render("paymentPlan/add_payment_plan", pageData)
	}

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok {
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	err = h.idempotencyService.UpdateResource(ctx, idempotencyKey, "member", id)
	if err != nil {
		return err
	}

	c.Status(fiber.StatusCreated)
	return h.renderPageByID(c, id)
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
	errorMap := req.validate()
	if len(errorMap) > 0 {
		paymentPlan, installments, err := h.service.Get(ctx, id)
		if err != nil {
			return err
		}
		res, err := mergetoResponse(*paymentPlan, req)
		if err != nil {
			return err
		}
		pageContext := page.PageContext{
			Mode:         "edit",
			UserAuthInfo: userAuthInfo,
			ActiveSection: "companies",
		}
		pageData := pageData{
			PaymentPlan:  res,
			Installments: installments,
			PageContext:  pageContext,
			Errors:       errorMap,
		}
		return c.Status(fiber.StatusBadRequest).Render("paymentPlan/payment_plan", pageData)
	}

	paymentPlan, err := toModel(req)
	if err != nil {
		return err
	}
	err = h.service.Update(ctx, id, paymentPlan)
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

	paymentPlan, _, err := h.service.Get(ctx, id)
	if err!=nil{
		return err
	}

	err = h.service.Cancel(ctx, id)
	if err!=nil{
		return err
	}

	return h.renderTableByID(c, paymentPlan.CompanyID)
}

func (h *PaymentPlanHandler) Restore(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	paymentPlan, _, err := h.service.Get(ctx, id)
	if err!=nil{
		return err
	}

	err = h.service.Restore(ctx, id)
	if err!=nil{
		return err
	}

	return h.renderTableByID(c, paymentPlan.CompanyID)
}

func (h *PaymentPlanHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	paymentPlan, _, err := h.service.Get(ctx, id)
	if err!=nil{
		return err
	}

	err = h.service.HardDelete(ctx, id)
	if err!=nil{
		return err
	}

	return h.renderTableByID(c, paymentPlan.CompanyID)
}



