package authports

import authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"

type ClaimsNormalizer interface {
	NormalizeClaims(any) *authdomain.AuthClaims
	DenormalizeClaims(authdomain.AuthClaims) any
}