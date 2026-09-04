package license

import (
	"errors"
	"fmt"

	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/gofiber/fiber/v2"
)

type LicenseHandler struct {
	service *LicenseService
}

func NewLicenseHandler(service *LicenseService) *LicenseHandler{
	return &LicenseHandler{service: service}
}

func (h *LicenseHandler) CheckUpdatedLicense(c *fiber.Ctx) error {
	valid, err := h.service.Check()
	if err != nil {
		return apperrors.NewInvalidLicenseError(fmt.Errorf("failed to check license: %w", err), "Error comprobando licencia, inténtelo nuevamente en unos segundos.")
	}
	h.service.SetAtomic(valid)
	if !valid {
		return apperrors.NewInvalidLicenseError(errors.New("license is not valid"), "")
	}
	return c.Redirect("/")
}