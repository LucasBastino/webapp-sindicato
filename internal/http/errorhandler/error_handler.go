package errorhandler

import (
	"errors"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/infra/logger"
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
	data := fiber.Map{"errorMsg": appErr.ClientMsg}
	if appErr.Type == "not_found" {
		data["errorTitle"] = "Registro no encontrado"
		if c.Get("HX-Request") == "" {
			data["goToDashboard"] = true
		}
	}

	if c.Get("HX-Request") == "true" {
		c.Set("HX-Retarget", "#app-modal-container")
		c.Set("HX-Reswap", "innerHTML")
		return c.Status(appErr.StatusCode).Render("error-modal", data)
	}

	return c.Status(appErr.StatusCode).Render("error/error_modal", data)
}

func renderLogin(c *fiber.Ctx, appErr *apperrors.AppError) error {
	// Session expired / invalid: always leave the current page and go to login.
	// HTMX must use HX-Redirect (full navigation); a 401 body would be swapped into
	// #app-modal-container by error-modal.js and leave a blank screen on close.
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", "/login")
		return c.SendStatus(appErr.StatusCode)
	}
	return c.Redirect("/login", fiber.StatusSeeOther)
}

func renderInvalidLicense(c *fiber.Ctx, appErr *apperrors.AppError) error {
	if c.Get("HX-Request") == "true" {
		c.Set("HX-Redirect", "/")
		return c.SendStatus(appErr.StatusCode)
	}
	return c.Status(appErr.StatusCode).Render("license/expired", fiber.Map{"errorMsg": appErr.ClientMsg})
}
