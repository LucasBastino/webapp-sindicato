package member

import (
	"errors"
	"fmt"
	"strconv"

	authports "github.com/LucasBastino/webapp-sindicato/internal/auth/ports"
	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/webapp-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/webapp-sindicato/internal/features/company"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/idempotency"
	"github.com/gofiber/fiber/v2"
)

type MemberHandler struct{
	service *MemberService
	companyService *company.CompanyService

	normalizer authports.ClaimsNormalizer
}

func NewMemberHandler(service *MemberService, companyService *company.CompanyService, normalizer authports.ClaimsNormalizer) *MemberHandler{
	return &MemberHandler{
		service: service,
		companyService: companyService,
		normalizer: normalizer,
	}
}

func (h *MemberHandler) RenderPage(c *fiber.Ctx) error {
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil{
		return err
	}
	return h.renderPageByID(c, id)
}

func (h *MemberHandler) renderPageByID(c *fiber.Ctx, id int) error {
	ctx := c.UserContext()
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	member, err := h.service.Get(ctx, id)
	if err!=nil {
		return err
	}
	res := toResponse(*member)
	pageContext	:= page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
	pageData := pageData{
		Member:      res,
		PageContext: pageContext,
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("member-content", pageData)
	}
	return c.Render("member/member", pageData)
}

func (h *MemberHandler) RenderTable(c *fiber.Ctx) error {
	
	switch c.Query("view"){
	case "electoral_list":
		return h.renderElectoralList(c)
	default:
		return h.renderDefaultTable(c)
	}

}

func (h *MemberHandler) RenderAddForm(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	res := response{}
	activeSection := "members"
	if raw := c.Params("company_id"); raw != "" {
		companyID, err := strconv.Atoi(raw)
		if err != nil {
			return apperrors.NewBadRequestError(fmt.Errorf("invalid company id: %w", err), "")
		}
		if companyID != 0 {
			res.CompanyID = companyID
			activeSection = "companies"
		}
	}
	pageData := pageData{
		Member:      res,
		PageContext: page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: activeSection},
	}
	return h.renderMemberPage(c, pageData, 0)
}

func (h *MemberHandler) renderDefaultTable(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var companyID *int
	companyIDParam := 0

	if raw := c.Params("company_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			return apperrors.NewBadRequestError(fmt.Errorf("invalid company id: %w", err), "")
		}
		if id != 0 {
			companyIDParam = id
			companyID = &companyIDParam
		}
	} else if raw := c.Query("company_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			return apperrors.NewBadRequestError(fmt.Errorf("invalid company id: %w", err), "")
		}
		if id != 0 {
			companyIDParam = id
			companyID = &companyIDParam
		}
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	searchKey := httpUtils.GetSearchKey(c)
	statusFilters := httpUtils.GetStatusFilters(c)

	filters := memberFilters{
		searchKey: searchKey,
		statuses:  statusFilters,
		companyID: companyID,
	}

	totalRows, err := h.service.Count(ctx, filters)
	if err != nil {
		return err
	}

	activeSection := "members"
	if companyIDParam != 0 {
		activeSection = "companies"
	}

	pageContext := page.PageContext{
		Mode:          "table",
		UserAuthInfo:  userAuthInfo,
		SearchKey:     searchKey,
		StatusFilters: statusFilters,
		ActiveSection: activeSection,
	}

	var members []tableResponse
	var emptyState page.EmptyState

	addHref := "/members/new"
	if companyIDParam != 0 {
		addHref = fmt.Sprintf("/companies/%d/members/new", companyIDParam)
	}

	if totalRows == 0 {
		if userAuthInfo.CanEdit("member") && statusFilters.ShowActive {
			emptyState = page.NewEmptyState("users", "afiliados", addHref, "Agregar afiliado")
		} else {
			emptyState = page.NewNoResultsEmptyState("users", "afiliados")
		}
	} else {
		currentPage := httpUtils.GetPageByQueryParam(c)
		pagination := page.BuildPagination(currentPage, totalRows)
		pageContext.Pagination = pagination

		memberList, err := h.service.List(ctx, filters, pagination.Offset)
		if err != nil {
			return err
		}
		members = toTableResponses(memberList)
	}

	tablePageData := tablePageData{
		Members:      members,
		TotalResults: totalRows,
		EmptyState:   emptyState,
		CompanyID:    companyIDParam,
		PageContext:  pageContext,
	}

	if companyIDParam != 0 {
		companyModel, err := h.companyService.Get(ctx, companyIDParam)
		if err != nil {
			return err
		}
		tablePageData.CompanyName = companyModel.Name
		number := ""
		if companyModel.CompanyNumber != nil {
			number = *companyModel.CompanyNumber
		}
		tablePageData.CompanyNav = page.NewCompanyNav(
			companyIDParam,
			companyModel.Name,
			"members",
			userAuthInfo.CanView("member"),
		).WithDetails(number, companyModel.Address, companyModel.Phone)
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("members-content", tablePageData)
	}
	return c.Render("member/members", tablePageData)
}

func (h *MemberHandler) renderElectoralList(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	members, err := h.service.GetElectoralList(ctx)
	if err != nil {
		return err
	}
	responses := toTableResponses(members)

	emptyState := page.EmptyState{}
	if len(responses) == 0 {
		emptyState = page.EmptyState{
			Icon:        "clipboard-list",
			Title:       "Padrón vacío",
			Description: "No hay afiliados habilitados para el padrón electoral.",
		}
	}

	tablePageData := tablePageData{
		Members:     responses,
		EmptyState:  emptyState,
		PageContext: page.PageContext{UserAuthInfo: userAuthInfo, ActiveSection: "reports"},
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("electoral-member-list-content", tablePageData)
	}
	return c.Render("member/electoral_member_list", tablePageData)
}


func (h *MemberHandler) Create(c *fiber.Ctx) error {
		record := c.Locals("idempotency_record")
	if record != nil {

		r, ok := record.(*idempotency.Record)
		if !ok{
			return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
		}

		c.Set("HX-Redirect", fmt.Sprintf("/members/%d", *r.ResourceID))
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
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse member req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		res, err := toResponseFromRequest(req)
		if err != nil {
			return err
		}
		activeSection := "members"
		if res.CompanyID != 0 {
			activeSection = "companies"
		}
		pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: activeSection}
		pageData := pageData{Member: res, PageContext: pageContext, Errors: errorMap}
		status := fiber.StatusBadRequest
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		return h.renderMemberPage(c, pageData, status)
	}


	// Si no tiene errores inserto el member en la DB y renderizo el su archivo
	member, err := toModel(req)
	if err!=nil{
		return err
	}

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok{
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	id, err := h.service.Create(ctx, member, idempotencyKey)
	if err!=nil{
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res, err := toResponseFromRequest(req)
			if err!=nil{
				return err
			}
			activeSection := "members"
			if res.CompanyID != 0 {
				activeSection = "companies"
			}
			PageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: activeSection}
			pageData := pageData{Member: res, PageContext: PageContext, Errors: errorMap}
			status := fiber.StatusConflict
			if c.Get("HX-Request") == "true" {
				status = fiber.StatusOK
			}
			return h.renderMemberPage(c, pageData, status)
		}
		return err
	}

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", fmt.Sprintf("/members/%d", id))
		return c.SendStatus(fiber.StatusOK)
	}
	return c.Redirect(fmt.Sprintf("/members/%d", id), fiber.StatusSeeOther)

}

func (h *MemberHandler) renderMemberPage(c *fiber.Ctx, pageData pageData, status int) error {
	if status != 0 {
		c.Status(status)
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("member-content", pageData)
	}
	return c.Render("member/member", pageData)
}

func (h *MemberHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	existing, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}
	if existing.DeletedAt != nil {
		return apperrors.NewBusinessError(errors.New("cannot update deleted member"), "No se puede editar un afiliado eliminado.")
	}
	
	var req request
	err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse member req body: %w", err), "")
	}


	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		res, err := mergetoResponse(*existing, req)
		if err != nil {
			return err
		}

		pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
		pageData := pageData{
			Member:      res,
			PageContext: pageContext,
			Errors:      errorMap,
		}
		status := fiber.StatusBadRequest
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		return h.renderMemberPage(c, pageData, status)
	}
	member, err := toModel(req)
	if err != nil {
		return err
	}
	err = h.service.Update(ctx, id, member)
	if err!=nil{
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res, err := mergetoResponse(*existing, req)
			if err != nil {
				return err
			}
			pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
			pageData := pageData{
				Member:      res,
				PageContext: pageContext,
				Errors:      errorMap,
			}
			status := fiber.StatusConflict
			if c.Get("HX-Request") == "true" {
				status = fiber.StatusOK
			}
			return h.renderMemberPage(c, pageData, status)
		}
		return err
	}
	
	return h.renderPageByID(c, id)

}

func (h *MemberHandler) SoftDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Obtengo el ID desde el path y lo elimino
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	err = h.service.SoftDelete(ctx, id)
	if err!=nil{
		return err
	}

	return h.renderDefaultTable(c)
}

func (h *MemberHandler) Restore(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	err = h.service.Restore(ctx, id)
	if err!=nil {
		return err
	}

	switch c.Get("view"){
	case "page":
		return h.renderPageByID(c, id)
	case "table":
		return h.renderDefaultTable(c)
	default:
		return apperrors.NewInternalError(errors.New("invalid query param while restoring member"), "")
	}
}

func (h *MemberHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// Obtengo el ID desde el path y lo elimino
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	err = h.service.HardDelete(ctx, id)
	if err!=nil{
		return err
	}

	return h.renderDefaultTable(c)
}




// func (h *CompanyHandler) RenderCompanyMembers(c *fiber.Ctx) error {
// 	ctx := c.UserContext()

// 	// obtengo la currentPage del path
// 	currentPage := httpUtils.GetPageByQueryParam(c)

// 	id, err := httpUtils.GetIDByParam(c, "company_id")
// 	if err!=nil {
// 		return err
// 	}

// 	includeInactive := c.Query("include-inactive") == "true"
// 	includeDeleted := c.Query("include-deleted") == "true"
// 	var searchKey string

// 	if c.Get("X-From-Delete") == "true" {
// 		// si estamos en deleteMode que el searchKey lo saque del header, ya que no se lo voy a mandar por el form
// 		// asi cuando elimino un miembro se quedan los miembros que busque antes menos el que elimine
// 		searchKey = c.Get("X-Search-Key")
// 	} else {
// 		// sino se lo mando por el form normalmente
// 		searchKey = c.FormValue("search-key")
// 	}

// 	// calculo la cantidad de resultados
// 	totalRows, err := h.service.Count(ctx, searchKey, includeInactive)
// 	if err!=nil {
// 		return err
// 	}

// 	if totalRows == 0 {
// 		// si no hay resultados renderizar esto
// 		return c.SendString(`<div class="no-result-file">No se encontraron afiliados</div>`)
// 	}
// 	// si hay resultados...

// 	// calcular totalPages
// 	totalPages, offset, someBefore, someAfter := httpUtils.GetpaginationMeta(currentPage, totalRows)
// 	limit := 10
// 	// busco los miembros y devuelvo el searchKey para usarlo nuevamente en la paginacion
// 	memberModels, err := h.service.ListMembers(ctx, id, limit, offset, searchKey, includeInactive, includeDeleted)
// 	if err!=nil {
// 		return err
// 	}
// 	memberTableResponses := member.toTableResponses(memberModels)

// 	// hago un array para poder recorrerlo y crear botones cuando hay menos de 10 paginas en el template
// 	totalPagesArray := httpUtils.GetTotalPagesArray(totalPages)
	
// 	// creo un map con todas las variables
// 	pagination := httpUtils.SetPagination(searchKey, currentPage, someBefore, someAfter, totalPages, totalPagesArray)
// 	companyMembersPage := CompanyMembersPageData{
// 		Members: memberTableResponses,
// 		CompanyID: id,
// 		TotalResults: totalRows,
// 		FromCompany: true,
// 		Pagination: pagination,
// 		UserAuthInfo: userAuthInfo,
// 		}

// 	return c.Render("member", companyMembersPage)
// }


// func (h *MemberHandler) validateForCreate(ctx context.Context, req request) map[string]string {
// 	errorMap := req.validate()
// 	if len(errorMap)>0{
// 		return errorMap
// 	}
// 	errorMap = h.service.ValidateInDBForCreate(ctx, req)
// 	if len(errorMap)>0{
// 		return errorMap
// 	}
// 	return nil
// }

// func (h *MemberHandler) validateForUpdate(ctx context.Context, id int, req request) map[string]string {
// 	errorMap := req.validate()
// 	if len(errorMap)>0{
// 		return errorMap
// 	}
// 	errorMap = h.service.ValidateInDBForUpdate(ctx, id, req)
// 	if len(errorMap)>0{
// 		return errorMap
// 	}
// 	return nil
// }




