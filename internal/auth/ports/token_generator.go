package authports

import (
	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
)

type TokenGenerator interface {
	Create(jwtSecret string, claims authdomain.AuthClaims) (string, error)
	Verify(jwtSecret string, token string) (*authdomain.AuthClaims, error)
}
