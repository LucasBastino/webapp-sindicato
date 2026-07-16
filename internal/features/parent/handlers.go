package parent

import (
	"errors"
	"fmt"

	authports "github.com/LucasBastino/app-sindicato/internal/auth/ports"
	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/gofiber/fiber/v2"
)

type ParentHandler struct{
	service *ParentService
	idempotencyService *idempotency.IdempotencyService

	normalizer authports.ClaimsNormalizer
}

func NewParentHandler(service *ParentService, idempotencyService *idempotency.IdempotencyService, normalizer authports.ClaimsNormalizer) *ParentHandler{
	return &ParentHandler{
		service: service,
		idempotencyService: idempotencyService,
		normalizer: normalizer,
	}
}

func (h *ParentHandler) RenderModal(c *fiber.Ctx) error {
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}
	return h.renderModalByID(c, id)
}

func (h *ParentHandler) renderModalByID(c *fiber.Ctx, id int) error {
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	
	parent, err := h.service.Get(ctx, id)
	if err!=nil {
		return err
	}
	res := toResponse(*parent)
	
	pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
	pageData := pageData{Parent: res, PageContext: pageContext}
	if c.Get("HX-Request") == "true" {
		return c.Render("parent-content", pageData)
	}
	return c.Render("parent/parent", pageData)
}

func (h *ParentHandler) RenderTable(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// calculo la cantidad de resultados
	memberID, err := httpUtils.GetIDByParam(c, "member_id")
	if err!=nil{
		return err
	}
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	totalRows, err := h.service.Count(ctx, memberID)
	if err!=nil {
		return err
	}

	pageContext := page.PageContext{
		Mode:          "edit",
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "members",
	}

	var parents []response
	var emptyState page.EmptyState

	if totalRows == 0 {
		emptyState = page.NewEmptyState(
			"users",
			"familiares",
			fmt.Sprintf("/members/%d/parents/new", memberID),
			"Agregar familiar",
		)
	} else {
		parentList, err := h.service.List(ctx, memberID)
		if err != nil {
			return err
		}
		parents = toTableResponses(parentList)
	}

	tablePageData := tablePageData{
		Parents:      parents,
		MemberID:     memberID,
		TotalResults: totalRows,
		EmptyState:   emptyState,
		PageContext:  pageContext,
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("parents-content", tablePageData)
	}
	return c.Render("parent/parents", tablePageData)
}

func (h *ParentHandler) RenderAddForm(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	pageContext := page.PageContext{Mode: "add", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
	pageData := pageData{Parent: response{}, PageContext: pageContext}
	if c.Get("HX-Request") == "true" {
		return c.Render("parent-content", pageData)
	}
	return c.Render("parent/parent", pageData)
}


func (h *ParentHandler) Create(c *fiber.Ctx) error {
		record := c.Locals("idempotency_record")
	if record != nil {

		r, ok := record.(*idempotency.Record)
		if !ok{
			return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
		}

		c.Set("HX-Redirect", fmt.Sprintf("/parents/%d", *r.ResourceID))
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
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse parent req body: %w", err), "")
	}
	
	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		res, err := toResponseFromRequest(req)
		if err!=nil{
			return err
		}
		pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "members"}
		pageData := pageData{Parent: res, PageContext: pageContext, Errors: errorMap}
		return c.Status(fiber.StatusBadRequest).Render("parent/parent", pageData)
	}
	parent, err := toModel(req)
	if err!=nil{
		return err
	}
	id, err := h.service.Create(ctx, parent)
	if err!=nil {
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res, err := toResponseFromRequest(req)
			if err!=nil{
				return err
			}
			pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "members"}
			pageData := pageData{Parent: res, PageContext: pageContext, Errors: errorMap}
			return c.Status(fiber.StatusConflict).Render("parent/parent", pageData)
		}
		return err
	}

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok{
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	err = h.idempotencyService.UpdateResource(ctx, idempotencyKey, "parent", id)
	if err != nil {
		return err
	}

	c.Status(fiber.StatusCreated)
	return h.renderModalByID(c, id)
}

func (h *ParentHandler) Update(c *fiber.Ctx) error {
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
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse parent req body: %w", err), "")
	}
	
	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		parent, err := h.service.Get(ctx, id)
		if err!=nil{
			return err
		}
		res, err := mergetoResponse(*parent, req)
		if err!=nil{
			return err
		}
		pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
		pageData := pageData{Parent: res, PageContext: pageContext, Errors: errorMap}
		return c.Status(fiber.StatusBadRequest).Render("parent/parent", pageData)
	}
	parent, err := toModel(req)
	if err!=nil{
		return err
	}
	err = h.service.Update(ctx, id, parent)
	if err!=nil {
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res, err := mergetoResponse(parent, req)
			if err!=nil{
				return err
			}
			pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
			pageData := pageData{Parent: res, PageContext: pageContext, Errors: errorMap}
			return c.Status(fiber.StatusConflict).Render("parent/parent", pageData)
		}
		return err
	}
	return h.renderModalByID(c, id)
}

func (h *ParentHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}
	err = h.service.HardDelete(ctx, id)
	if err!=nil{
		return err
	}
	return h.RenderTable(c)
}


/* func (h *ParentHandler) DeleteParentFromMember(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "parent_id")
	if err!=nil {
		return err
	}
	err = h.service.HardDelete(ctx, id)
	if err!=nil{
		return err
	}
	return h.RenderParentTable(c)
} */



