package payment

import (
	"fmt"
	"strconv"

	authports "github.com/LucasBastino/app-sindicato/internal/auth/ports"
	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct{
	service *PaymentService

	normalizer authports.ClaimsNormalizer
}

func NewPaymentHandler(service *PaymentService, normalizer authports.ClaimsNormalizer) *PaymentHandler{
	return &PaymentHandler{
		service: service,
		normalizer: normalizer,
	}
}

// func (h *PaymentHandler) validate(req request) map[string]string{
// 	errorMap := req.ValidateBasic()
// 	if len(errorMap) > 0 {
// 		return errorMap
// 	}
// 	return nil
// }

// func (h *PaymentHandler) RenderAddForm(c *fiber.Ctx) error {
// 	companyID, err := httpUtils.GetIDByParam(c, "payment_id")
// 	if err!=nil {
// 		return err
// 	}

// 	model := Payment{CompanyID: companyID}
// 	data := fiber.Map{"payment": model, "mode": "add"}
// 	return c.Render("payment_file", data)
// }

func (h *PaymentHandler) RenderModal(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err!=nil {
		return err
	}

	payment, err := h.service.Get(ctx, id)
	if err != nil{
		return err
	}

	res := toResponse(*payment)

	pageContext := page.PageContext{
		Mode: "edit",
		UserAuthInfo: userAuthInfo,
		ActiveSection: "companies",
	}
	pageData := pageData{
		Payment: res,
		CompanyID: companyID,
		PageContext: pageContext,
	}
	return c.Render("payment/payment", pageData)
}

func (h *PaymentHandler) RenderGrid(c *fiber.Ctx) error {
	ctx := c.UserContext()
	year := c.Query("year")

	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err!=nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	companyName := c.Get("X-Company-Name")

	totalRows, err := h.service.Count(ctx, companyID)
	if err!=nil{
		return err
	}

	renderPayments := func(data gridPageData) error {
		if c.Get("HX-Request") == "true" {
			return c.Render("payments-content", data)
		}
		return c.Render("payment/payments", data)
	}

	if totalRows == 0 {
		pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "edit", ActiveSection: "companies"}
		tablePageData := gridPageData{CompanyID: companyID, CompanyName: companyName, PageContext: pageContext}
		return renderPayments(tablePageData)
	}

	years, err := h.service.ListPaymentYears(ctx, companyID)
	if err!=nil{
		return err
	}

	pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "edit", ActiveSection: "companies"}
	tablePageData := gridPageData{CompanyID: companyID, CompanyName: companyName, Years: years, PageContext: pageContext}

	if year == "0" {
		payments, err := h.service.List(ctx, companyID, years[0])
		if err != nil {
			return err
		}
		responses := toGridResponses(payments)
		tablePageData.Payments = responses
		tablePageData.Year = years[0]
		return renderPayments(tablePageData)
	} else {
		yearInt, err := strconv.Atoi(year)
		if err!=nil{
			return apperrors.NewInternalError(fmt.Errorf("failed to parse year: %w", err), "")
		}
		payments, err := h.service.List(ctx, companyID, yearInt)
		if err != nil {
			return err
		}
		responses := toGridResponses(payments)
		tablePageData.Payments = responses
		tablePageData.Year = yearInt
		return renderPayments(tablePageData)
	}
}

func (h *PaymentHandler) CreatePayments(c *fiber.Ctx) error{
	ctx := c.UserContext()
	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err!=nil {
		return err
	}
	return h.service.CreateRemainingPaymentsByID(ctx, nil, companyID)
}

func (h *PaymentHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	var req request
	err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		payment, err := h.service.Get(ctx, id)
		if err!=nil{
			return err
		}
		res, err := mergetoResponse(*payment, req)
		if err!=nil{
			return err
		}
		pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo}
		pageData := pageData{Payment: res, PageContext: pageContext, Errors: errorMap}
		return c.Status(fiber.StatusBadRequest).Render("payment_file", pageData)
	}
	payment, err := toModel(req)
	if err!=nil{
		return err
	}
	err = h.service.Update(ctx, id, payment)
	if err!=nil{
		return err
	}

/* 	paymentRes := PaymentModelToRes(modelFromDB)

	pageContext := PageContext{Mode: "edit", ResourceRoles:httpUtils.GetUserPermissions(c, h.normalizer)}
	pageData := pageData{Payment: paymentRes, CompanyID: companyID, PageContext: pageContext}
	return c.Render("paymentFile", pageData) */
	return c.SendStatus(fiber.StatusNoContent)
}



// func (h *PaymentHandler) Create(c *fiber.Ctx) error {
// 	ctx := c.UserContext()

// 	var req request
// 	err := c.BodyParser(&req)
// 	if err!=nil{
// 		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse payment req body: %w", err), "")
// 	}
	
// 	companyID, err := httpUtils.GetIDByParam(c, "payment_id")
// 	if err!=nil {
// 		return err
// 	}

// 	companyName := c.Get("X-Company-Name")

// 	errorMap := h.validate(req)
// 	if len(errorMap) > 0 {
// 		res := toResponseFromRequest(req)
// 		pageContext := page.PageContext{Mode: "add"}
// 		pageData := pageData{Payment: res, CompanyID: companyID, CompanyName: companyName, PageContext: pageContext}
// 		return c.Render("payment_file", pageData)
// 	}
// 	model := toModel(req)
// 	id, err := h.service.Create(ctx, model)
// 	if err!=nil {
// 		return err
// 	}

// 	/* pageContext := PageContext{Mode: "edit", ResourceRoles:httpUtils.GetUserPermissions(c, h.normalizer)}
// 	pageData := pageData{Payment: paymentRes, CompanyID: companyID, PageContext: pageContext}
// 	return c.Render("paymentFile", pageData) */
// 	path := fmt.Sprintf("/payment/%d/file", id)
// 	return c.Status(fiber.StatusCreated).Render("redirect", fiber.Map{"path": path})
// }

