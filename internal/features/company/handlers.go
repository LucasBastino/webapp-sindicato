package company

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

type CompanyHandler struct{
	service *CompanyService
	idempotencyService *idempotency.IdempotencyService

	normalizer authports.ClaimsNormalizer
}

func NewCompanyHandler(service *CompanyService, idempotencyService *idempotency.IdempotencyService, normalizer authports.ClaimsNormalizer) *CompanyHandler{
	return &CompanyHandler{
		service: service,
		idempotencyService: idempotencyService,
		normalizer: normalizer,
	}
}


func (h *CompanyHandler) RenderPage(c *fiber.Ctx) error {
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil{
		return err
	}
	return h.renderPageByID(c, id)
}

func (h *CompanyHandler) renderPageByID(c *fiber.Ctx, id int) error{
	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	company, err := h.service.Get(ctx, id)
	if err!=nil {
		return err
	}

	res := toResponse(*company)
	
	pageContext := page.PageContext{UserAuthInfo: userAuthInfo, Mode: "edit", ActiveSection: "companies"}
	pageData := pageData{Company: res, PageContext: pageContext}
	pageData.WithPaymentTable = false

	if c.Get("HX-Request") == "true" {
		return c.Render("company-content", pageData)
	}
	return c.Render("company/company", pageData)
}

func (h *CompanyHandler) RenderTable(c *fiber.Ctx) error {

	switch c.Query("view"){
	case "options":
		return h.renderOptions(c)
	default:
		return h.renderDefaultTable(c)
	}	
}

func (h *CompanyHandler) renderDefaultTable(c *fiber.Ctx) error{
	ctx := c.UserContext()
	
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	searchKey := httpUtils.GetSearchKey(c)
	statusFilters := httpUtils.GetStatusFilters(c)

	filters := companyFilters{
		searchKey: searchKey,
		statuses: statusFilters,
	}

	// calculo la cantidad de resultados
	totalRows, err := h.service.Count(ctx, filters)
	if err!=nil {
		return err
	}

	pageContext := page.PageContext{
		UserAuthInfo:  userAuthInfo,
		SearchKey:     searchKey,
		StatusFilters: statusFilters,
		ActiveSection: "companies",
	}

	var companies []tableResponse
	var emptyState page.EmptyState

	if totalRows == 0 {
		emptyState = page.NewEmptyState("building-2", "empresas", "/companies/new", "Agregar empresa")
	} else {
		currentPage := httpUtils.GetPageByQueryParam(c)
		pagination := page.BuildPagination(currentPage, totalRows)
		pageContext.Pagination = pagination

		companyList, err := h.service.List(ctx, filters, pagination.Offset)
		if err != nil {
			return err
		}
		companies = toTableResponses(companyList)
	}

	pageData := tablePageData{
		Companies:    companies,
		TotalResults: totalRows,
		EmptyState:   emptyState,
		PageContext:  pageContext,
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("companies-content", pageData)
	}
	return c.Render("company/companies", pageData)
}

func (h *CompanyHandler) renderOptions(c *fiber.Ctx) error {
	ctx := c.UserContext()
	// calculo la cantidad de resultados
	searchKey := c.FormValue("search-key")
	companyOptions, err := h.service.ListForSelect(ctx, searchKey)
	if err!=nil {
		return err
	}
	
	if len(companyOptions) == 0 {
		// si no hay resultados renderizar esto
		return c.SendString(`<div class="no-result">No se encontraron empresas</div>`)
	}

	// si hay resultados...
	optionResponses := toOptionResponses(companyOptions)
	pageData := optionsPageData{Companies: optionResponses}

	return c.Render("companyTableSelect", pageData)
	
}

func (h *CompanyHandler) RenderAddForm(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}
	// le paso un company vacio para que los campos del form aparezcan vacios
	pageData := pageData{
		Company: response{},
		PageContext: page.PageContext{UserAuthInfo: userAuthInfo, Mode: "add", ActiveSection: "companies"},
	}
	return c.Render("company/company", pageData)
}


func (h *CompanyHandler) Create(c *fiber.Ctx) error {

	record := c.Locals("idempotency_record")
	if record != nil {

		r, ok := record.(*idempotency.Record)
		if !ok{
			return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
		}

		c.Set("HX-Redirect", fmt.Sprintf("/companies/%d", *r.ResourceID))
		return c.SendStatus(200)

	}

	ctx := c.UserContext()

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	// creo la request(dto) (no tiene timestamps ni id) y le copio los datos provenientes del form
	var req request
	err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse company req body: %w", err), "")
	}
	
	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		res := toResponseFromRequest(req)
		
		pageContext := page.PageContext{
			UserAuthInfo: userAuthInfo,
			Mode:         "create",
			ActiveSection: "companies",
		}
		pageData := pageData{
			Company:     res,
			PageContext: pageContext,
			Errors:      errorMap,
		}
		return c.Status(fiber.StatusBadRequest).Render("company/company", pageData)
	}


	// aca lo paso de companyReq a model puro
	company := toModel(req)

	// lo inserto en la DB y lo retorno con mas datos (id, timestamps)
	id, err := h.service.Create(ctx, company)
	if err!=nil{
		// chequeo duplicados
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res := toResponseFromRequest(req)

			pageContext := page.PageContext{
				UserAuthInfo: userAuthInfo,
				Mode:         "create",
				ActiveSection: "companies",
			}
			pageData := pageData{
				Company:     res,
				PageContext: pageContext,
				Errors:      errorMap,
			}
			
			return c.Status(fiber.StatusConflict).Render("company/company", pageData)
		}
		return err
	}	
	

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok{
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}
	
	err = h.idempotencyService.UpdateResource(ctx, idempotencyKey, "company", id)
	if err != nil {
		return err
	}

	/* // lo paso a res
	companyRes := modelToRes(modelFromDB)

	pageContext := PageContext{UserAuthInfo:getUserPermissions(c, h.normalizer), Mode: "edit"}
	pageData := pageData{Company: companyRes,	NumberOfMembers: 0,	PageContext: pageContext}

	return c.Render("companyFile", pageData) */
	c.Status(fiber.StatusCreated)
	return h.renderPageByID(c, id)
}

func (h *CompanyHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	// obtengo el id del param
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
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse company req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		company, err := h.service.Get(ctx, id)
		if err!=nil{
			return err
		}
		res := mergetoResponse(*company, req)
		
		pageContext := page.PageContext{
			UserAuthInfo: userAuthInfo,
			Mode:         "edit",
			ActiveSection: "companies",
		}
		pageData := pageData{
			Company:     res,
			PageContext: pageContext,
			Errors:      errorMap,
		}

		return c.Status(fiber.StatusBadRequest).Render("company/company", pageData)
	}

	company := toModel(req)

	err = h.service.Update(ctx, id, company)
	if err!=nil{
		mapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			res := toResponseFromRequest(req)
			pageContext := page.PageContext{
				UserAuthInfo: userAuthInfo,
				Mode:         "edit",
				ActiveSection: "companies",
			}
			pageData := pageData{
				Company:     res,
				PageContext: pageContext,
				Errors:      errorMap,
			}

			return c.Status(fiber.StatusConflict).Render("company/company", pageData)
		}
		return err
	}	
	
	return h.renderPageByID(c, id)
}

func (h *CompanyHandler) SoftDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	err = h.service.SoftDelete(ctx, id)
	if err!=nil {
		return err
	}

	return h.renderDefaultTable(c)
}

func (h *CompanyHandler) Restore(c *fiber.Ctx) error {
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
		return apperrors.NewInternalError(errors.New("invalid query param while restoring company"), "")
	}
}

func (h *CompanyHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	err = h.service.HardDelete(ctx, id)
	if err!=nil {
		return err
	}

	return h.renderDefaultTable(c)
}






// func (h *CompanyHandler) validateForCreate(ctx context.Context, req request) map[string]string{
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

// func (h *CompanyHandler) validateForUpdate(ctx context.Context, id int, req request) map[string]string{
// 	errorMap := req.validate()
// 	if len(errorMap)>0{
// 		return errorMap
// 	}

// 	// errorMap = h.service.ValidateInDBForUpdate(ctx, id, req)
// 	// if len(errorMap)>0{
// 	// 	return errorMap
// 	// }
// 	return nil
// }