package errorhandler

import (
	"errors"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
)

type ErrorHandler struct{
	logger logger.Logger
}


func HandleError(c *fiber.Ctx, err error) error {

	var appErr *apperrors.AppError

	if !errors.As(err, &appErr){
		appErr = apperrors.NewInternalError(err, "")
	}

	switch appErr.RenderType {

	case apperrors.RenderTypeToast:
		return renderToast(c, appErr)
	
	case apperrors.RenderTypeModal:
		return renderModal(c, appErr)
	
	case apperrors.RenderTypeLogin:
		return renderLogin(c, appErr)
	
	case apperrors.RenderTypeInvalidLicense:
		return renderInvalidLicense(c, appErr)
	}

	return renderModal(c, appErr)
}

func renderToast(c *fiber.Ctx, appErr *apperrors.AppError) error {
	return c.Status(appErr.StatusCode).Render("toast", fiber.Map{"errorMsg": appErr.ClientMsg})
}

func renderModal(c *fiber.Ctx, appErr *apperrors.AppError) error {
	return c.Status(appErr.StatusCode).Render("error_modal", fiber.Map{"errorMsg": appErr.ClientMsg})
}

func renderLogin(c *fiber.Ctx, appErr *apperrors.AppError) error {
	return c.Status(appErr.StatusCode).Render("user/login", fiber.Map{
		"Username": "",
		"Errors":   map[string]string{},
	})
}

func renderInvalidLicense(c *fiber.Ctx, appErr *apperrors.AppError) error {
	return c.Status(appErr.StatusCode).Render("invalid_license", fiber.Map{"errorMsg": appErr.ClientMsg})
}