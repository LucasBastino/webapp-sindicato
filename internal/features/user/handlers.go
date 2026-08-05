package user

import (
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) canChangePassword(authAdmin bool, authUserID, targetID int) bool {
	return authAdmin || authUserID == targetID
}

func requireCurrentPassword(actorAdmin bool, actorID, targetID int, targetAdmin bool) bool {
	if !actorAdmin {
		return true // no-admin solo toca la propia
	}
	return actorID == targetID || targetAdmin
}

func (h *UserHandler) RenderPanel(c *fiber.Ctx) error {
	ctx := c.UserContext()
	users, err := h.service.List(ctx)
	if err != nil {
		return err
	}
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	data := TablePageData{
		Users: ToTableResponses(users),
		PageContext: page.PageContext{
			UserAuthInfo:  userAuthInfo,
			ActiveSection: "users",
		},
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("users-content", data)
	}
	return c.Render("user/users", data)
}

func (h *UserHandler) RenderAddForm(c *fiber.Ctx) error {
	userAuthInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	data := PageData{
		PageContext: page.PageContext{
			UserAuthInfo:  userAuthInfo,
			ActiveSection: "users",
		},
	}
	if c.Get("HX-Request") == "true" {
		return c.Render("register-content", data)
	}
	return c.Render("user/register", data)
}

func (h *UserHandler) RenderChangePasswordModal(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	authInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	if !h.canChangePassword(authInfo.Admin, authInfo.UserID, userID) {
		return apperrors.NewForbiddenError(
			fmt.Errorf("user %d cannot change password for user %d", authInfo.UserID, userID),
			"No tenés permisos para cambiar esta contraseña.",
		)
	}
	target, err := h.service.Get(ctx, userID)
	if err != nil {
		return err
	}
	return c.Render("change-password-modal", PasswordModalData{
		ID:                     userID,
		RequireCurrentPassword: requireCurrentPassword(authInfo.Admin, authInfo.UserID, userID, target.Admin),
	})
}

func (h *UserHandler) RenderUpdatePermissionsModal(c *fiber.Ctx) error {
	ctx := c.UserContext()
	userID, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	authInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	if userID == authInfo.UserID {
		return apperrors.NewBusinessError(apperrors.ErrCannotEditOwnPermissions, "No podés editar tus propios permisos.")
	}
	user, err := h.service.Get(ctx, userID)
	if err != nil {
		return err
	}
	return c.Render("update-permissions-modal", PermissionsModalData{
		ID:            user.ID,
		Username:      user.Username,
		Admin:         user.Admin,
		ResourceRoles: user.ResourceRoles,
	})
}

func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	authInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}
	if !h.canChangePassword(authInfo.Admin, authInfo.UserID, id) {
		return apperrors.NewForbiddenError(
			fmt.Errorf("user %d cannot change password for user %d", authInfo.UserID, id),
			"No tenés permisos para cambiar esta contraseña.",
		)
	}

	target, err := h.service.Get(ctx, id)
	if err != nil {
		return err
	}
	requireCurrent := requireCurrentPassword(authInfo.Admin, authInfo.UserID, id, target.Admin)

	var req passwordRequest
	err = c.BodyParser(&req)
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse change password req body: %w", err), "")
	}

	errorMap := req.Validate(requireCurrent)
	if len(errorMap) > 0 {
		c.Set("HX-Retarget", "#app-modal-container")
		c.Set("HX-Reswap", "innerHTML")
		status := fiber.StatusBadRequest
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		return c.Status(status).Render("change-password-modal", PasswordModalData{
			ID:                     id,
			RequireCurrentPassword: requireCurrent,
			CurrentPassword:        req.CurrentPassword,
			Password:               req.Password,
			ConfirmPassword:        req.ConfirmPassword,
			Errors:                 errorMap,
		})
	}

	err = h.service.ChangePassword(ctx, id, req, requireCurrent)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidCurrentPassword) {
			errorMap = map[string]string{"current_password": "La contraseña actual no es correcta."}
		} else if errors.Is(err, apperrors.ErrInvalidNewPassword) {
			errorMap = map[string]string{"password": "Debes ingresar una contraseña distinta a la actual."}
		} else {
			return err
		}
		c.Set("HX-Retarget", "#app-modal-container")
		c.Set("HX-Reswap", "innerHTML")
		status := fiber.StatusBadRequest
		if c.Get("HX-Request") == "true" {
			status = fiber.StatusOK
		}
		return c.Status(status).Render("change-password-modal", PasswordModalData{
			ID:                     id,
			RequireCurrentPassword: requireCurrent,
			CurrentPassword:        req.CurrentPassword,
			Password:               req.Password,
			ConfirmPassword:        req.ConfirmPassword,
			Errors:                 errorMap,
		})
	}

	c.Set("HX-Retarget", "#app-modal-container")
	c.Set("HX-Reswap", "innerHTML")
	return c.SendString("")
}

func (h *UserHandler) UpdatePermissions(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	authInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	var req permissionsRequest
	err = c.BodyParser(&req)
	if err != nil {
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse user permissions req body: %w", err), "")
	}

	errorMap := req.Validate()
	if len(errorMap) > 0 {
		c.Set("HX-Retarget", "#users-modal-container")
		c.Set("HX-Reswap", "innerHTML")
		admin, roles := req.ToPermissions()
		return c.Status(fiber.StatusBadRequest).Render("update-permissions-modal", PermissionsModalData{
			ID:            id,
			Admin:         admin,
			ResourceRoles: roles,
			Errors:        errorMap,
		})
	}

	admin, permissions := req.ToPermissions()

	err = h.service.UpdatePermissions(ctx, id, authInfo.UserID, admin, permissions)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPermissions) {
			errorMap = map[string]string{"resourceRoles": "Un admin debe tener permisos de editor en todos los recursos."}
			c.Set("HX-Retarget", "#users-modal-container")
			c.Set("HX-Reswap", "innerHTML")
			return c.Status(fiber.StatusBadRequest).Render("update-permissions-modal", PermissionsModalData{
				ID:            id,
				Admin:         admin,
				ResourceRoles: permissions,
				Errors:        errorMap,
			})
		}
		return err
	}

	return h.RenderPanel(c)
}

func (h *UserHandler) HardDelete(c *fiber.Ctx) error {
	ctx := c.UserContext()
	id, err := httpUtils.GetIDByParam(c, "id")
	if err != nil {
		return err
	}

	authInfo, err := httpUtils.GetUserAuthInfo(c)
	if err != nil {
		return err
	}

	err = h.service.HardDelete(ctx, id, authInfo.UserID)
	if err != nil {
		return err
	}

	return h.RenderPanel(c)
}
