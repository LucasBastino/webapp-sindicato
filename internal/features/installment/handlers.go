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

func (h InstallmentHandler) buildPageData(c *fiber.Ctx, installment Installment, errors map[string]string) (pageData, error) {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return pageData{}, err
	}

	planStatus, err := h.service.GetPaymentPlanStatus(c.UserContext(), installment.PaymentPlanID)
	if err != nil {
		return pageData{}, err
	}
	planCancelled := planStatus == "cancelled"
	canEdit := userAuthInfo.CanEdit("company") && !planCancelled

	return pageData{
		Installment:        toResponse(installment),
		PageContext:        page.PageContext{UserAuthInfo: userAuthInfo, ActiveSection: "companies"},
		Errors:             errors,
		CanEditInstallment: canEdit,
		PlanCancelled:      planCancelled,
	}, nil
}

func (h InstallmentHandler) RenderModal(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	installment, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}

	pageData, err := h.buildPageData(c, *installment, nil)
	if err != nil {
		return err
	}

	if c.Get("HX-Request") == "true" {
		return c.Render("installment-modal", pageData)
	}
	return c.Render("installment/installment", pageData)
}

func (h *InstallmentHandler) Update(c *fiber.Ctx) error {
	ctx := c.UserContext()

	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	var req request
	err = c.BodyParser(&req)
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse installment req body: %w", err), "")
	}

	req.trim()
	errorMap := req.validate()
	if len(errorMap) > 0 {
		installment, err := h.service.Get(ctx, id)
		if err != nil {
			return err
		}
		res, err := mergetoResponse(*installment, req)
		if err != nil {
			return err
		}

		pageData, err := h.buildPageData(c, *installment, errorMap)
		if err != nil {
			return err
		}
		pageData.Installment = res
		if c.Get("HX-Request") == "true" {
			return c.Status(fiber.StatusOK).Render("installment-modal", pageData)
		}
		return c.Status(fiber.StatusBadRequest).Render("installment/installment", pageData)
	}

	installment, err := toModel(req)
	if err != nil {
		return err
	}
	err = h.service.Update(ctx, id, installment)
	if err != nil {
		return err
	}

	c.Set("HX-Retarget", "#app-modal-container")
	c.Set("HX-Reswap", "innerHTML")
	c.Set("HX-Trigger", "paymentPlanUpdated")
	return c.SendString("")
}
