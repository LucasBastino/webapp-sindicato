package installment

import (
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type InstallmentHandler struct {
	service *InstallmentService
}

func NewInstallmentHandler(service *InstallmentService) *InstallmentHandler {
	return &InstallmentHandler{
		service: service,
	}
}



func (h InstallmentHandler) RenderModal(c *fiber.Ctx) error {
	ctx := c.UserContext()
	
	id, err := httpUtils.GetIDByParam(c, "id")
	if err!=nil {
		return err
	}

	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err!=nil{
		return err
	}

	installment, err := h.service.Get(ctx, id)
	if err != nil{
		return err
	}

	res := toResponse(*installment)

	pageContext := page.PageContext{UserAuthInfo: userAuthInfo, ActiveSection: "companies"}
	pageData := pageData{
		Installment: res,
		PageContext: pageContext,
	}

	return c.Render("installment/installment", pageData)
}


func (h *InstallmentHandler) Update(c *fiber.Ctx) error {
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
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse installment req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		installment, err := h.service.Get(ctx, id)
		if err!=nil{
			return err
		}
		res, err := mergetoResponse(*installment, req)
		if err!=nil{
			return err
		}

		pageContext := page.PageContext{UserAuthInfo: userAuthInfo, ActiveSection: "companies"}
		pageData := pageData{
			Installment: res,
			PageContext: pageContext,
			Errors:      errorMap,
		}
		return c.Status(fiber.StatusBadRequest).Render("installment/installment", pageData)
	}
	
	installment, err := toModel(req)
	if err!=nil{
		return err
	}
	err = h.service.Update(ctx, id, installment)
	if err!=nil{
		return err
	}

	return h.RenderModal(c)
}

