package user

import (
	"errors"
	"fmt"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/app-sindicato/internal/common/page"
	httpUtils "github.com/LucasBastino/app-sindicato/internal/common/utils/http"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{
	service *UserService

}

func NewUserHandler(service *UserService) *UserHandler{
	return &UserHandler{
		service: service,
	}
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
    userID, err := c.ParamsInt("user_id")
    if err != nil {
        return err
    }
    return c.Render("changePasswordModal", fiber.Map{"ID": userID})
}

func (h *UserHandler) RenderUpdatePermissionsModal(c *fiber.Ctx) error {
    userID, err := c.ParamsInt("user_id")
    if err != nil {
        return err
    }
    return c.Render("update_permissions_modal", fiber.Map{"ID": userID})
}


func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
    ctx := c.UserContext()
    id, err := c.ParamsInt("id")
    if err != nil {
        return err
    }

    var req passwordRequest
    err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse change password req body: %w", err), "")
	}

    errorMap := req.Validate()
    if len(errorMap) > 0 {
        // c.Set("HX-Retarget", "#modal-container")
        return c.Status(fiber.StatusBadRequest).Render("change_password_modal", errorMap)
    }

    err = h.service.ChangePassword(ctx, id, req)
    if err != nil {
        if errors.Is(err, apperrors.ErrInvalidCurrentPassword){
            errorMap["current-password"] =  "La contraseña actual no es correcta."
        } else{
            // si es un error del hasher, retornar el error
            return err
        }
        // c.Set("HX-Retarget", "#modal-container")
        res := ToPasswordResponseFromRequest(req)

        return c.Render("change_password_modal", res)
    }

    // c.Set("HX-Retarget", "#modal-container")
    return c.SendString("")
}

func (h *UserHandler) UpdatePermissions(c *fiber.Ctx) error {
    ctx := c.UserContext()
    id, err := c.ParamsInt("id")
    if err != nil {
        return err
    }

    var req permissionsRequest
	err = c.BodyParser(&req)
	if err!=nil{
		return apperrors.NewBadRequestError(fmt.Errorf("failed to parse user permissions req body: %w", err), "")
	}

    errorMap := req.Validate()
    if len(errorMap) > 0 {
        // c.Set("HX-Retarget", "#modal-container")
        return c.Status(fiber.StatusBadRequest).Render("update_permissions_modal", errorMap)
    }

    admin, permissions := req.ToPermissions()
    
    err = h.service.UpdatePermissions(ctx, id, admin, permissions)
    if err != nil {
        if errors.Is(err, apperrors.ErrInvalidPermissions){
            errorMap["resourceRoles"] = "Un admin debe tener permisos de editor en todos los recursos."
        } else{
            return err
        }
        // c.Set("HX-Retarget", "#modal-container")
        res := ToPermissionsResponseFromRequest(req)

        return c.Render("update_permissions_modal", res)
    }

    // c.Set("HX-Retarget", "#modal-container")
    return c.SendString("")


}

func (h *UserHandler) HardDelete(c *fiber.Ctx) error {
    ctx := c.UserContext()
    // Obtengo el ID desde el path y lo elimino
    id, err := httpUtils.GetIDByParam(c, "id")
    if err!=nil {
        return err
    }

    err = h.service.HardDelete(ctx, id)
    if err!=nil{
        return err
    }

    return h.RenderPanel(c)
}
