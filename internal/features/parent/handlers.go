package parent

import (
	"errors"
	"fmt"

	authports "github.com/LucasBastino/webapp-sindicato/internal/auth/ports"
	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/webapp-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/webapp-sindicato/internal/features/member"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
	"github.com/gofiber/fiber/v2"
)

type ParentHandler struct{
	service *ParentService
	memberService *member.MemberService

	normalizer authports.ClaimsNormalizer
}

func NewParentHandler(service *ParentService, memberService *member.MemberService, normalizer authports.ClaimsNormalizer) *ParentHandler{
	return &ParentHandler{
		service: service,
		memberService: memberService,
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
	pageData := pageData{Parent: res, MemberID: parent.MemberID, PageContext: pageContext}
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
		if userAuthInfo.CanEdit("member") {
			emptyState = page.NewEmptyState(
				"users",
				"familiares",
				fmt.Sprintf("/members/%d/parents/new", memberID),
				"Agregar familiar",
			)
		} else {
			emptyState = page.NewNoResultsEmptyState("users", "familiares")
		}
	} else {
		parentList, err := h.service.List(ctx, memberID)
		if err != nil {
			return err
		}
		parents = toTableResponses(parentList)
	}

	memberModel, err := h.memberService.Get(ctx, memberID)
	if err != nil {
		return err
	}
	memberName := fmt.Sprintf("%s, %s", memberModel.LastName, memberModel.Name)
	memberNumber := ""
	if memberModel.MemberNumber != nil {
		memberNumber = *memberModel.MemberNumber
	}

	tablePageData := tablePageData{
		Parents:      parents,
		MemberID:     memberID,
		MemberName:   memberName,
		Member: memberCardData{
			ID:           memberModel.ID,
			Name:         memberModel.Name,
			LastName:     memberModel.LastName,
			MemberNumber: memberNumber,
			Dni:          memberModel.Dni,
			Phone:        memberModel.Phone,
			CompanyID:    memberModel.CompanyID,
			CompanyName:  memberModel.CompanyName,
		},
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
	memberID, err := httpUtils.GetIDByParam(c, "member_id")
	if err != nil {
		return err
	}
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	pageContext := page.PageContext{Mode: "add", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
	pageData := pageData{Parent: response{}, MemberID: memberID, PageContext: pageContext}
	if c.Get("HX-Request") == "true" {
		return c.Render("parent-content", pageData)
	}
	return c.Render("parent/parent", pageData)
}


func (h *ParentHandler) Create(c *fiber.Ctx) error {
	memberID, err := httpUtils.GetIDByParam(c, "member_id")
	if err != nil {
		return err
	}

	record := c.Locals("idempotency_record")
	if record != nil {
		r, ok := record.(*idempotency.Record)
		if !ok {
			return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
		}
		_ = r
		c.Set("HX-Redirect", fmt.Sprintf("/members/%d/parents", memberID))
		return c.SendStatus(200)
	}

	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	var req request
	err = c.BodyParser(&req)
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse parent req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()

	renderAddError := func(errors map[string]string, status int) error {
		res, err := toResponseFromRequest(req)
		if err != nil {
			return err
		}
		pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "members"}
		pd := pageData{Parent: res, MemberID: memberID, PageContext: pageContext, Errors: errors}
		if c.Get("HX-Request") == "true" {
			return c.Status(fiber.StatusOK).Render("parent-content", pd)
		}
		return c.Status(status).Render("parent/parent", pd)
	}

	if len(errorMap) > 0 {
		return renderAddError(errorMap, fiber.StatusBadRequest)
	}

	parent, err := toModel(req)
	if err != nil {
		return err
	}
	parent.MemberID = memberID

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok {
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	_, err = h.service.Create(ctx, parent, idempotencyKey)
	if err != nil {
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			return renderAddError(errorMap, fiber.StatusConflict)
		}
		return err
	}

	c.Set("HX-Redirect", fmt.Sprintf("/members/%d/parents", memberID))
	return c.SendStatus(fiber.StatusCreated)
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

	renderError := func(res response, memberID int, errors map[string]string) error {
		pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
		pd := pageData{Parent: res, MemberID: memberID, PageContext: pageContext, Errors: errors}
		if c.Get("HX-Request") == "true" {
			return c.Render("parent-content", pd)
		}
		return c.Render("parent/parent", pd)
	}

	if len(errorMap) > 0 {
		existing, err := h.service.Get(ctx, id)
		if err != nil {
			return err
		}
		res, err := mergetoResponse(*existing, req)
		if err != nil {
			return err
		}
		return renderError(res, existing.MemberID, errorMap)
	}

	parent, err := toModel(req)
	if err != nil {
		return err
	}
	err = h.service.Update(ctx, id, parent)
	if err != nil {
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			existing, err := h.service.Get(ctx, id)
			if err != nil {
				return err
			}
			res, err := mergetoResponse(*existing, req)
			if err != nil {
				return err
			}
			return renderError(res, existing.MemberID, errorMap)
		}
		return err
	}
	return h.renderModalByID(c, id)
}

func (h *ParentHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}
	parent, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}
	memberID := parent.MemberID
	err = h.service.HardDelete(ctx, id)
	if err != nil {
		return err
	}
	c.Set("HX-Redirect", fmt.Sprintf("/members/%d/parents", memberID))
	return c.SendStatus(fiber.StatusOK)
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



