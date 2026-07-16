package authadapters

import (
	"errors"
	"fmt"

	authdomain "github.com/LucasBastino/app-sindicato/internal/auth/domain"
	authports "github.com/LucasBastino/app-sindicato/internal/auth/ports"
	"github.com/LucasBastino/app-sindicato/internal/common/apperrors"
	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenGenerator struct{
	normalizer authports.ClaimsNormalizer
}

func NewJwtTokenGenerator(normalizer authports.ClaimsNormalizer) *JwtTokenGenerator{
	return &JwtTokenGenerator{normalizer: normalizer}
}

func (j *JwtTokenGenerator) Create(JWTSecret string, defaulthClaims authdomain.AuthClaims) (string, error){
	// hago el type assertion defaulthClaims.(authClaims.JwtClaims) para que go confie que es un authClaims.JwtClaims
	// jwtClaims := j.normalizer.DenormalizeClaims(defaulthClaims).(authClaims.JwtClaims)
	jwtMapClaims := j.normalizer.DenormalizeClaims(defaulthClaims).(jwt.MapClaims)
	// JwtClaims struct contiene la variable MapClaims que es un jwt.MapClaims que implementa la interfaz jwt.Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtMapClaims)
	signedToken, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", apperrors.NewUnauthorizedError(fmt.Errorf("failed to sign token: %w", err), "")
	}
	return signedToken, nil
}


func (j *JwtTokenGenerator) Verify(JWTSecret string, tokenStr string) (*authdomain.AuthClaims, error){
	// jwt.Parse devuelve un jwt.Claims, no un jwt.MapClaims
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return nil, apperrors.NewUnauthorizedError(fmt.Errorf("failed to parse token: %w", err), "")
	}
	// jwt.Parse no siempre devuelve error cuando el token esta expirado o mal formado, a veces devuelve un token no valido, por lo tanto hay que verificarlo
	if !token.Valid {
    return nil, apperrors.NewUnauthorizedError(errors.New("invalid token"), "")
}
	// aca creo que va la funcion JwtClaimsToAuthClaims, pero capaz no hace falta y con lo de abajo ya esta, probar
	// porque en create si o si necesitas pasar un jwtmapclaims, y no se puede hacer type assertion porque el jwtmapclaims no tiene getExpiry por ejemplo
	// pero si se puede pasar de jwtmapclaims a authclaims porque tiene todos los campos? no se, me parece raro
	// claims := token.Claims.(authClaims.AuthClaims)
	jwtMapclaims, ok := token.Claims.(jwt.MapClaims)
	if !ok{
		return nil, apperrors.NewUnauthorizedError(errors.New("invalid token claims type"), "")
	}
	defaulthClaims := j.normalizer.NormalizeClaims(jwtMapclaims)
	if defaulthClaims.Exp.IsZero() || defaulthClaims.Sub == 0 || defaulthClaims.ResourceRoles == nil {
		return nil, apperrors.NewUnauthorizedError(errors.New("invalid token claims data"), "")
	}
	return defaulthClaims, nil
}

/* 
// este temp sirve solamente para la conversion del jwt a authclaim
type TempAuthClaims struct{
	userId string
	Admin bool
	permissions []string
	exp time.Time
}

// y le adjudico estos metodos solamente para que cumpla con la interfaz AuthClaims
func (t TempAuthClaims) GetUserID() string {return t.userId}
func (t TempAuthClaims) GetAdmin() bool {return t.Admin}
func (t TempAuthClaims) GetUserPermissions() []string {return t.permissions}
func (t TempAuthClaims) GetExpiry() time.Time {return t.exp}

func AuthClaimsToJwtClaims(claims authClaims.AuthClaims) jwt.MapClaims {
	return jwt.MapClaims{
		"userId": claims.GetUserID(),
		"Admin": claims.GetAdmin(),
		"permissions": claims.GetUserPermissions(),
		"exp": claims.GetExpiry(), 
	}
}

func JwtClaimsToAuthClaims(claims jwt.MapClaims) TempAuthClaims{
	return TempAuthClaims{
		userId: claims["userId"].(string),
		Admin: claims["Admin"].(bool),
		permissions: claims["permissions"].([]string),
		exp: time.Unix(int64(claims["exp"].(float64)), 0),
	}
} */