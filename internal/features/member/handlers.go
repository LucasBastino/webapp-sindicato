package member

import (
	"errors"
	"fmt"
	"strconv"

	authports "github.com/LucasBastino/app-sindicato/internal/auth/ports"
	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/gofiber/fiber/v2"
)

type MemberHandler struct{
	service *MemberService
	idempotencyService *idempotency.IdempotencyService

	normalizer authports.ClaimsNormalizer
}

func NewMemberHandler(service *MemberService, idempotencyService *idempotency.IdempotencyService, normalizer authports.ClaimsNormalizer) *MemberHandler{
	return &MemberHandler{
		service: service,
		idempotencyService: idempotencyService,
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
	pageData := pageData{Member: res, PageContext: pageContext}
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
	companyID, err := httpUtils.GetIDByParam(c, "company_id")
	if err!=nil{
		return err
	}
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	res := response{}
	
	if companyID != 0{
		res.CompanyID = companyID
	}
	
	pageData := pageData{Member: res, PageContext: page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "members"}}
	return c.Render("member/member", pageData)
}

func (h *MemberHandler) renderDefaultTable(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var companyID *int
	var companyIDParam int
	var err error

	raw := c.Query("id")
	if raw != "" {
		// si el param existe obtengo el companyID de la URL
		companyIDParam, err = strconv.Atoi(raw)
		if err != nil {
			return apperrors.NewBadRequestError(fmt.Errorf("invalid company id: %w", err), "")
		}
		companyID = &companyIDParam
	}

	// si el companyID es 0 el valor tiene que ser nil
	if companyIDParam != 0{
		companyID = &companyIDParam
	} else{
		companyID = nil
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	// FILTROS
	searchKey := httpUtils.GetSearchKey(c)
	
	// obtengo los filtros del form
	statusFilters := httpUtils.GetStatusFilters(c)

	// declaro los filtros
	filters := memberFilters{
		searchKey: searchKey,
		statuses: statusFilters,
		companyID: companyID,
	}

	// cuento la cantidad de resultados obtenidos con estos filtros
	totalRows, err := h.service.Count(ctx, filters)
	if err!=nil {
		return err
	}

	pageContext := page.PageContext{
		Mode:          "table",
		UserAuthInfo:  userAuthInfo,
		SearchKey:     searchKey,
		StatusFilters: statusFilters,
		ActiveSection: "members",
	}

	var members []tableResponse
	var emptyState page.EmptyState

	if totalRows == 0 {
		emptyState = page.NewEmptyState("users", "afiliados", "/members/new", "Agregar afiliado")
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
		PageContext:  pageContext,
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("members-content", tablePageData)
	}
	return c.Render("member/members", tablePageData)
}

func (h *MemberHandler) renderElectoralList(c *fiber.Ctx) error {
	ctx := c.UserContext()
	members, err:= h.service.GetElectoralList(ctx)
	if err!=nil{
		return err
	}
	responses := toTableResponses(members)

	return c.Render("electoral_member_list", fiber.Map{"members": responses})
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
		if err!=nil{
			return err
		}		
		PageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "members"}
		pageData := pageData{Member: res, PageContext: PageContext, Errors: errorMap}
		return c.Status(fiber.StatusBadRequest).Render("member/member", pageData)
	}


	// Si no tiene errores inserto el member en la DB y renderizo el su archivo
	member, err := toModel(req)
	if err!=nil{
		return err
	}
	id, err := h.service.Create(ctx, member)
	if err!=nil{
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res, err := toResponseFromRequest(req)
			if err!=nil{
				return err
			}
			PageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "members"}
			pageData := pageData{Member: res, PageContext: PageContext, Errors: errorMap}
			return c.Status(fiber.StatusConflict).Render("member/member", pageData)
		}
		return err
	}

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok{
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	err = h.idempotencyService.UpdateResource(ctx, idempotencyKey, "member", id)
	if err != nil {
		return err
	}

	c.Status(fiber.StatusCreated)
	return h.renderPageByID(c, id)

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
	
	var req request
	err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse member req body: %w", err), "")
	}


	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		member, err := h.service.Get(ctx, id)
		if err != nil {
			return err
		}
		res, err := mergetoResponse(*member, req)
		if err != nil {
			return err
		}

		pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
		pageData := pageData{Member: res, PageContext: pageContext, Errors: errorMap}
		return c.Status(fiber.StatusBadRequest).Render("member/member", pageData)
	}
	member, err := toModel(req)
	if err != nil {
		return err
	}
	err = h.service.Update(ctx, id, member)
	if err!=nil{
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res, err := mergetoResponse(member, req)
			if err != nil {
				return err
			}
			pageContext := page.PageContext{Mode: "edit", UserAuthInfo: userAuthInfo, ActiveSection: "members"}
			pageData := pageData{Member: res, PageContext: pageContext, Errors: errorMap}
			return c.Status(fiber.StatusConflict).Render("member/member", pageData)
		}
		return err
	}
	
	//  hacer un redirect mejor 

	/* data := fiber.Map{"member": m, "mode": "edit", "companies": companies, "companyName": companyName, "createdAt": createdAt, "updatedAt": updatedAt}
	data["canDelete"] = c.Locals("claims").(jwt.MapClaims)["canDelete"]
	data["canWrite"] = c.Locals("claims").(jwt.MapClaims)["canWrite"]
	return c.Render("member_file", data) */
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




