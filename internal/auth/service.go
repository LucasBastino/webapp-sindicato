package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	authdomain "github.com/LucasBastino/webapp-sindicato/internal/auth/domain"
	authports "github.com/LucasBastino/webapp-sindicato/internal/auth/ports"
	"github.com/LucasBastino/webapp-sindicato/internal/common/apperrors"
	"github.com/LucasBastino/webapp-sindicato/internal/config"
	"github.com/LucasBastino/webapp-sindicato/internal/features/user"
	"github.com/LucasBastino/webapp-sindicato/internal/security/password"
)

type AuthService struct {
	repo *AuthRepository
	userService *user.UserService

	passwordHasher password.Hasher
	tokenGenerator authports.TokenGenerator
    
    cfg config.AuthConfig
}

// Precomputed bcrypt (cost 14) used when the username is missing so login
// timing stays close to the "user exists, wrong password" path.
var loginDummyPasswordHash = []byte("$2a$14$Sy2SjaW9pdiZ2kcP7EJSwOU3t8bssnq0mLhsgfepIhNopzUizHw62")

func NewAuthService(repo *AuthRepository, userService *user.UserService, passwordHasher password.Hasher, tokenGenerator authports.TokenGenerator, cfg config.AuthConfig ) *AuthService{
	return &AuthService{
		repo: repo,
        userService: userService,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
        cfg: cfg,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, error){
    user, err := s.userService.GetByUsername(ctx, username)
	if err!=nil {
		return "", "", apperrors.NewDatabaseError(err, "")
	}

    if user == nil{
		_ = s.passwordHasher.Compare(loginDummyPasswordHash, password)
		return "", "", fmt.Errorf("%w", apperrors.ErrInvalidLogin)
	}

	err = s.passwordHasher.Compare([]byte(user.PasswordHash), password)
    if err != nil {
        return "", "", fmt.Errorf("%w", apperrors.ErrInvalidLogin)
	}

	claims := authdomain.AuthClaims{
		Sub:           user.ID,
		Exp:           time.Now().Add(s.cfg.AccessTokenTTL),
		Iat:           time.Now(),
		Username:      user.Username,
		Admin:         user.Admin,
		ResourceRoles: user.ResourceRoles,
	}
    
    refreshToken, err := s.CreateRefreshToken(ctx, user.ID)
    if err != nil {
        return "", "", apperrors.NewInternalError(err, "")
    }

	accessToken, err := s.tokenGenerator.Create(s.cfg.JWTSecret, claims)
	if err!=nil {
		return "", "", apperrors.NewInternalError(err, "")
	}

    return refreshToken, accessToken, nil
}

// func (s *AuthService) GenerateHash(password string) ([]byte, error){
//     return s.passwordHasher.Generate(password)
// }

// func (s *AuthService) CompareHashAndPassword(hash []byte, password string) error{
// 	return s.passwordHasher.Compare(hash, password)
// }

// func (s *AuthService) CreateToken(claims authdomain.AuthClaims) (string, error){
// 	return s.tokenGenerator.Create(claims)
// }

func (s *AuthService) VerifyAccessToken(token string) (*authdomain.AuthClaims, error) {
    return s.tokenGenerator.Verify(s.cfg.JWTSecret, token)
}

func (s *AuthService) CreateRefreshToken(ctx context.Context, userID int) (string, error) {
    // generás un string random seguro
    rawToken, err := generateRandomToken()
    if err != nil {
        return "", err
    }

    // hasheás con SHA-256 para guardar en DB
    hash := hashToken(rawToken)

    expiresAt := time.Now().Add(s.cfg.RefreshTokenTTL) // 7 días

    err = s.repo.SaveRefreshToken(ctx, userID, hash, expiresAt)
    if err!=nil{
		return "", apperrors.NewDatabaseError(err, "")
	}

    return rawToken, nil
}

func (s *AuthService) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	err := s.repo.RevokeRefreshToken(ctx, tokenHash)
    if err!=nil{
		return apperrors.NewDatabaseError(err, "")
	}

    return nil
}

func (s *AuthService) RevokeAllSessionsForUser(ctx context.Context, userID int) error {
	if err := s.repo.RevokeAllRefreshTokensForUser(ctx, userID); err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	return nil
}

func (s *AuthService) RevokeOtherSessionsForUser(ctx context.Context, userID int, rawRefreshToken string) error {
	exceptHash := hashToken(rawRefreshToken)
	if err := s.repo.RevokeAllRefreshTokensForUserExcept(ctx, userID, exceptHash); err != nil {
		return apperrors.NewDatabaseError(err, "")
	}
	return nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, rawRefreshToken string) (string, string, *authdomain.AuthClaims, error) {
    hash := hashToken(rawRefreshToken)

    refreshToken, err := s.repo.GetRefreshToken(ctx, hash)
    if err != nil {
        return "", "", nil, apperrors.NewDatabaseError(fmt.Errorf("failed to fetch refresh token: %w", err), "")
    }

    if refreshToken == nil{
        return "", "", nil, apperrors.NewUnauthorizedError(errors.New("refresh token not found"), "")
    }

    if refreshToken.RevokedAt != nil {
		// Possible reuse of a rotated token: invalidate all sessions for that user.
		_ = s.repo.RevokeAllRefreshTokensForUser(ctx, refreshToken.UserID)
        return "", "", nil, apperrors.NewUnauthorizedError(errors.New("refresh token revoked"), "Sesión inválida.")
    }
    if refreshToken.ExpiresAt.Before(time.Now()) {
        return "", "", nil, apperrors.NewUnauthorizedError(errors.New("refresh token expired"), "Sesión expirada.")
    }

    user, err := s.userService.Get(ctx, refreshToken.UserID)
    if err != nil {
        return "", "", nil, err
    }

	if err := s.repo.RevokeRefreshToken(ctx, hash); err != nil {
		return "", "", nil, apperrors.NewDatabaseError(err, "")
	}

	newRefreshToken, err := s.CreateRefreshToken(ctx, refreshToken.UserID)
	if err != nil {
		return "", "", nil, err
	}

    claims := authdomain.AuthClaims{
        Sub:            refreshToken.UserID,
        Exp:            time.Now().Add(s.cfg.AccessTokenTTL),
        Iat:            time.Now(),
        Username:       user.Username,
        Admin:          user.Admin,
        ResourceRoles:  user.ResourceRoles,
    }

    accessToken, err := s.tokenGenerator.Create(s.cfg.JWTSecret, claims)
    if err != nil {
        return "", "", nil, apperrors.NewInternalError(err, "")
    }

    return accessToken, newRefreshToken, &claims, nil
}

func (s *AuthService) Register(ctx context.Context, user user.User, password string, idempotencyKey string) (int, error) {
    byteHash, err := s.passwordHasher.Generate(password)
	if err != nil {
		return 0, apperrors.NewInternalError(fmt.Errorf("failed to generate password hash: %w", err), "")
	}
    passWordHash := string(byteHash)
    user.PasswordHash = passWordHash
    
    return s.userService.Create(ctx, user, idempotencyKey)
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {
    if rawRefreshToken == "" {
        return nil
    }
    tokenHash := hashToken(rawRefreshToken)
    err := s.repo.RevokeRefreshToken(ctx, tokenHash)
    if err != nil {
        return apperrors.NewDatabaseError(err, "")
    }
    return nil
}

// func (s *AuthService) GetUserByUsername(ctx context.Context, username string) (user.User, error) {
//     return s.userService.GetByUsername(ctx, username)
// }

