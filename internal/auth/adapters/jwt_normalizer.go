package authadapters

import (
	"time"

	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JwtNormalizer struct {
}

func NewJwtNormalizer() *JwtNormalizer{
	return &JwtNormalizer{}
}

func (n JwtNormalizer) NormalizeClaims(jwtMapClaims any) *authdomain.AuthClaims{

	// despues en otro lugar puedo verificar si defaultclaims esta vacio hacer tal cosa
	// o si userID o is admin o ... esta vacio, hacer tal cosa
	claims, ok := jwtMapClaims.(jwt.MapClaims)
	if !ok{
		return nil
	}

	var sub int
	subVal, ok := claims["sub"].(float64)
	if !ok {
		sub = 0
	} else {
		sub = int(subVal)
	}


	var exp time.Time
	// primero se hace type assertion a float64 porque los numeros json se interpretan con ese tipo
	expVal, ok := claims["exp"].(float64)
	if !ok{
		exp = time.Time{}
	} else {
		//  despues el float64 lo convertis a int64 y el int64 a time.Time usando Unix
		exp = time.Unix(int64(expVal), 0)
	}

	var iat time.Time
	iatVal, ok := claims["iat"].(float64)
	if !ok{
		iat = time.Time{}
	} else {
		iat = time.Unix(int64(iatVal), 0)
	}

	username, ok := claims["username"].(string)
	if !ok {
		username = ""
	}

	admin, ok := claims["admin"].(bool)
	if !ok {
		admin = false
	}
	resourceRoles, ok := claims["resourceRoles"].(map[string]any)
	if !ok {
		resourceRoles = nil
	}

	return &authdomain.AuthClaims{
		Sub:           sub,
		Exp:           exp,
		Iat:           iat,
		Username:      username,
		Admin:         admin,
		ResourceRoles: resourceRoles,
	}
}

func (n JwtNormalizer) DenormalizeClaims(defaultClaims authdomain.AuthClaims) any {
	return jwt.MapClaims{
		// hay que ponerle .Unix() para convertir el tipo time.Time a int64
		"sub":           defaultClaims.Sub,
		"exp":           defaultClaims.Exp.Unix(),
		"iat":           defaultClaims.Iat.Unix(),
		"username":      defaultClaims.Username,
		"admin":         defaultClaims.Admin,
		"resourceRoles": defaultClaims.ResourceRoles,
	}
}
