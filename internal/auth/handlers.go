package auth

import (
	"errors"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/LucasBastino/app-sindicato/internal/features/user"
	"github.com/LucasBastino/app-sindicato/internal/infra/logger"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *AuthService

	logger logger.Logger
}

func NewAuthHandler(service *AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
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

func (h *AuthHandler) renderRegisterForm(c *fiber.Ctx, pageData user.PageData, status int) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	pageData.PageContext = page.PageContext{
		UserAuthInfo:  userAuthInfo,
		ActiveSection: "users",
	}
	if c.Get("HX-Request") == "true" {
		status = fiber.StatusOK
		return c.Status(status).Render("register-content", pageData)
	}
	return c.Status(status).Render("user/register", pageData)
}

func (h *AuthHandler) renderUsersPanel(c *fiber.Ctx) error {
	ctx := c.UserContext()
	users, err := h.service.userService.List(ctx)
	if err != nil {
		return err
	}
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	pageData := user.TablePageData{
		Users: user.ToTableResponses(users),
		PageContext: page.PageContext{
			UserAuthInfo:  userAuthInfo,
			ActiveSection: "users",
		},
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("users-content", pageData)
	}
	return c.Render("user/users", pageData)
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	record := c.Locals("idempotency_record")
	if record != nil {
		c.Set("HX-Redirect", "/user_panel")
		return c.SendStatus(200)
	}

	ctx := c.UserContext()

	var req user.Request
	if err := c.BodyParser(&req); err != nil {
		return apperrors.NewBadRequestError(err, "")
	}

	errorMap := req.Validate()
	if len(errorMap) > 0 {
		return h.renderRegisterForm(c, user.PageData{
			User:   user.ToResponseFromRequest(req),
			Errors: errorMap,
		}, fiber.StatusBadRequest)
	}

	userModel := user.ToModel(req)

	idempotencyKey, ok := c.Locals("idempotency_key").(string)
	if !ok {
		return apperrors.NewInternalError(errors.New("invalid idempotency record type in context"), "")
	}

	_, err := h.service.Register(ctx, userModel, req.Password, idempotencyKey)
	if err != nil {
		errorMap = map[string]string{}
		if errors.Is(err, apperrors.ErrInvalidPermissions) {
			errorMap["resourceRoles"] = "Un admin debe tener permisos de editor en todos los recursos."
		}
		user.MapDBDuplicateError(err, errorMap)
		if len(errorMap) > 0 {
			return h.renderRegisterForm(c, user.PageData{
				User:   user.ToResponseFromRequest(req),
				Errors: errorMap,
			}, fiber.StatusBadRequest)
		}
		return err
	}

	c.Set("HX-Push-Url", "/user_panel")
	return h.renderUsersPanel(c)
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
	}

	clearCookies(c)

	return c.Redirect("/login")
}
