package backup

import (
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
)

type BackUpHandler struct {
	service *BackUpService
	logger logger.Logger
}

func NewBackUpHandler(service *BackUpService, logger logger.Logger) *BackUpHandler {
	return &BackUpHandler{
		service: service,
		logger: logger,
	}
}

func (h BackUpHandler) Backup(c *fiber.Ctx) error {
	err := h.service.BackUp()
	if err!=nil{
		return err
	}
	h.logger.Info("Backup succesfull")
	return c.Render("backup_modal", fiber.Map{"done": true})
}