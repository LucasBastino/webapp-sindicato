package auth

import (
	"errors"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/features/user"
	"github.com/LucasBastino/app-sindicato/internal/infra/idempotency"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct{
	service *AuthService
	idempotencyService *idempotency.IdempotencyService

	logger logger.Logger
}

func NewAuthHandler(service *AuthService, idempotencyService *idempotency.IdempotencyService, logger logger.Logger) *AuthHandler{
	return &AuthHandler{
		service: service,
		idempotencyService: idempotencyService,
		logger: logger,
	}
}

func (ah *AuthHandler) RenderInsufficientPermissions(c *fiber.Ctx) error {
	return c.Render("insufficientPermissions", fiber.Map{})
}

func (ah *AuthHandler) RenderExpiredSession(c *fiber.Ctx) error {
	return c.Render("expiredSession", fiber.Map{})
}

func (h *AuthHandler) RenderLogin(c *fiber.Ctx) error {
	return c.Render("user/login", LoginPageData{Errors: map[string]string{}})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	record := c.Locals("idempotency_record")
	if record != nil {
		c.Set("HX-Redirect", "/users")
		return c.SendStatus(200)
	}

	ctx := c.UserContext()

	var req user.Request 
	c.BodyParser(&req)

	errorMap := req.Validate()
	if len(errorMap) > 0 {
		res := user.ToResponseFromRequest(req)
		pageData := user.PageData{
			User:   res,
			Errors: errorMap,
		}
		return c.Status(fiber.StatusBadRequest).Render("user/register", pageData)
	}

	userModel := user.ToModel(req)

	id, err := h.service.Register(ctx, userModel, req.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPermissions){
			errorMap["resourceRoles"] = "Un admin debe tener permisos de editor en todos los recursos."
		}
		user.MapDBDuplicateError(err, errorMap)
	}

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok{
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	err = h.idempotencyService.UpdateResource(ctx, idempotencyKey, "user", id)
	if err != nil {
		return err
	}

	users, err := h.service.userService.List(ctx)
	if err!=nil{
		return err
	}

	res := user.ToTableResponses(users)

	pageData := user.TablePageData{
		Users:  res,
		Errors: errorMap,
	}

	return c.Render("user/users", pageData)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx := c.UserContext()

	username, password := c.FormValue("user"), c.FormValue("password")

	refreshToken, accessToken, err := h.service.Login(ctx, username, password)
	if err != nil {
		errorMap := map[string]string{}
		if errors.Is(err, apperrors.ErrInvalidLoginUser) {
			errorMap["user"] = "Usuario incorrecto."
		} else if errors.Is(err, apperrors.ErrInvalidLoginPassword) {
			errorMap["password"] = "Contraseña incorrecta."
		} else {
			return err
		}
		return c.Status(fiber.StatusUnauthorized).Render("user/login", LoginPageData{
			Username: username,
			Errors:   errorMap,
		})
	}

	refreshCookie := createRefreshCookie(refreshToken, h.service.cfg.RefreshTokenTTL)
	c.Cookie(&refreshCookie)

	accessCookie := createAccessCookie(accessToken, h.service.cfg.AccessTokenTTL)
	c.Cookie(&accessCookie)

	return c.Redirect("/dashboard")
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx := c.UserContext()
	refreshToken := c.Cookies("refresh_token")
    
    err := h.service.Logout(ctx, refreshToken)
    if err != nil {
        h.logger.Error("logout continued after refresh token revoke failed", "err", err)
		// no retorno nada, dejo que se desloguee con clearCookies
    }
	
	clearCookies(c)

	return c.Redirect("/login")

}