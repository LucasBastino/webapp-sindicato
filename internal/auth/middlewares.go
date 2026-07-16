package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	userauthinfo "github.com/LucasBastino/app-sindicato/internal/features/user/authinfo"
	"github.com/LucasBastino/app-sindicato/internal/license"
	"github.com/gofiber/fiber/v2"
)

type AuthMiddleware struct{
	authService *AuthService
	licenseService *license.LicenseService
}

func NewAuthMiddleware(authService *AuthService, licenseService *license.LicenseService) *AuthMiddleware{
	return &AuthMiddleware{
		authService: authService,
		licenseService: licenseService,
	}
}


func (m *AuthMiddleware) VerifyToken(c *fiber.Ctx) error{

	// whitelist
	if c.Path() == "/favicon.ico" || strings.HasPrefix(c.Path(), "/static/") {
   		return c.Next()
	}

	if c.Path() == "/login" {
   		return c.Next()
	}

	token := c.Cookies("access_token")
	claims, err := m.authService.VerifyAccessToken(token)
	if err != nil {
		// token ausente, malformado o expirado → intentar con refresh token
		return m.handleRefresh(c)
	}
	if claims == nil {
		return apperrors.NewUnauthorizedError(errors.New("invalid token claims data"), "")
	}
	// si el token expiró
	if claims.Exp.Before(time.Now()) {
		return m.handleRefresh(c)
	}
	userAuthInfo := userauthinfo.UserAuthInfo{
		UserID: claims.Sub,
        Admin:   claims.Admin,
        ResourceRoles: claims.ResourceRoles,
    }
	// envio los claims y userAuthInfo a locals y dejo pasar al siguiente middleware
	c.Locals("claims", claims)
	c.Locals("userAuthInfo", userAuthInfo)
	return c.Next()
}

func (m *AuthMiddleware) handleRefresh(c *fiber.Ctx) error {
    ctx := c.UserContext()
    refreshToken := c.Cookies("refresh_token")
    if refreshToken == "" {
		clearCookies(c)
        return apperrors.NewUnauthorizedError(errors.New("failed to get refresh_token cookie from context"), "")
    }

    // hasheas el token y buscas en DB
    newAccessToken, claims, err := m.authService.RefreshAccessToken(ctx, refreshToken)
    if err != nil {
		clearCookies(c)
        return apperrors.NewUnauthorizedError(fmt.Errorf("failed to refresh access token: %w", err), "")
    }

    // seteás la nueva cookie
    accessCookie := createAccessCookie(newAccessToken, m.authService.cfg.AccessTokenTTL)
    c.Cookie(&accessCookie)

	userAuthInfo := userauthinfo.UserAuthInfo{
		UserID: claims.Sub,
        Admin:   claims.Admin,
        ResourceRoles: claims.ResourceRoles,
    }

    c.Locals("claims", claims)
    c.Locals("userAuthInfo", userAuthInfo)
    return c.Next()
}

func (m *AuthMiddleware) VerifyRole(resource string, requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error{
		userAuthInfo, ok := c.Locals("userAuthInfo").(userauthinfo.UserAuthInfo)
		if !ok{
			return apperrors.NewInternalError(errors.New("user auth info missing from context"), "")
		}
		if userAuthInfo.Admin{
			return c.Next()
		}
		userRole := userAuthInfo.ResourceRoles[resource]
		if userRole == requiredRole{
			return c.Next()
		}
		roleLevels := map[interface{}]int{
			"editor": 2,
			"viewer": 1,
		}
		if roleLevels[userRole] > roleLevels[requiredRole]{
			return c.Next()
		}
		return apperrors.NewForbiddenError(fmt.Errorf("user_id=%d lacks role %q on resource %q (has %q)", userAuthInfo.UserID, requiredRole, resource, userRole), "No tenés permisos para realizar esta acción.")
	}
}

func (m *AuthMiddleware) VerifyAdmin(c *fiber.Ctx) error{
	userAuthInfo, ok := c.Locals("userAuthInfo").(userauthinfo.UserAuthInfo)
	if !ok{
		return apperrors.NewInternalError(errors.New("user auth info missing from context"), "")
	}
	if userAuthInfo.Admin{
		return c.Next()
	}
	return apperrors.NewForbiddenError(fmt.Errorf("user %d is not admin", userAuthInfo.UserID), "Esta acción requiere permisos de administrador.")
}



func (m *AuthMiddleware) VerifyAtomicLicense (c *fiber.Ctx) error {
	// whitelist para verifyLicense
	if c.Path() == "/verifyLicense" {
   		return c.Next()
	}
	
	valid := m.licenseService.IsValid()
	if !valid {
		fmt.Println("licencia falsa")
		clearCookies(c)
		return apperrors.NewInvalidLicenseError(errors.New("license is not valid"), "")
	}

	// fmt.Println("paso por el middleware, es true")
	return c.Next()
}