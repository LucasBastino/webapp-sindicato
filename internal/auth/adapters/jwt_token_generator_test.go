package authadapters

import (
	"testing"
	"time"

	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
	"github.com/golang-jwt/jwt/v5"
)

func TestJwtVerify_AcceptsHS256(t *testing.T) {
	gen := NewJwtTokenGenerator(NewJwtNormalizer())
	secret := "0123456789abcdef0123456789abcdef"
	claims := authdomain.AuthClaims{
		Sub:           1,
		Username:      "admin",
		Admin:         true,
		ResourceRoles: map[string]any{"company": "editor"},
		Exp:           time.Now().Add(time.Hour),
		Iat:           time.Now(),
	}

	token, err := gen.Create(secret, claims)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := gen.Verify(secret, token)
	if err != nil {
		t.Fatalf("Verify HS256: %v", err)
	}
	if got.Sub != 1 || got.Username != "admin" {
		t.Fatalf("unexpected claims: %+v", got)
	}
}

func TestJwtVerify_RejectsUnexpectedAlg(t *testing.T) {
	gen := NewJwtTokenGenerator(NewJwtNormalizer())
	secret := "0123456789abcdef0123456789abcdef"

	mapClaims := jwt.MapClaims{
		"sub":           float64(1),
		"username":      "admin",
		"admin":         true,
		"resourceRoles": map[string]any{"company": "editor"},
		"exp":           float64(time.Now().Add(time.Hour).Unix()),
		"iat":           float64(time.Now().Unix()),
	}

	tokenHS384 := jwt.NewWithClaims(jwt.SigningMethodHS384, mapClaims)
	signedHS384, err := tokenHS384.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign HS384: %v", err)
	}
	if _, err := gen.Verify(secret, signedHS384); err == nil {
		t.Fatal("expected Verify to reject HS384")
	}

	tokenNone := jwt.NewWithClaims(jwt.SigningMethodNone, mapClaims)
	signedNone, err := tokenNone.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}
	if _, err := gen.Verify(secret, signedNone); err == nil {
		t.Fatal("expected Verify to reject alg=none")
	}
}
