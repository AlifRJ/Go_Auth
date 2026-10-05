package service

import (
	"context"
	"errors"
	"time"

	"github.com/AlifRJ/Go_Auth/app/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid identity or password")
	ErrInvalidToken       = errors.New("invalid or expired refresh token")
)

type AuthService struct {
	userRepo model.UserRepository
	authRepo model.AuthRepository
	accessTokenAuth  *jwtauth.JWTAuth
	refreshTokenAuth *jwtauth.JWTAuth
}

func NewAuthService(userRepo model.UserRepository, authRepo model.AuthRepository, accessAuth *jwtauth.JWTAuth, refreshAuth *jwtauth.JWTAuth,) *AuthService {
	return &AuthService{
		userRepo: userRepo, 
		authRepo: authRepo, 
		accessTokenAuth:  accessAuth,
		refreshTokenAuth: refreshAuth,
	}
}

func (s *AuthService) Login(ctx context.Context, identity, password string) (string, string, *model.User, error) {
	// Search User in repository
	user, err := s.userRepo.GetByEmailOrUsername(ctx, identity)
	if err != nil || user == nil {
		return "", "", nil, ErrInvalidCredentials
	}

	// Verify Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", nil, ErrInvalidCredentials
	}

	user.Password = ""

	// Generate Access & Refresh Token
	accessToken, refreshToken, err := s.generateTokenPair(ctx, user.ID, user.Name, user.Username, user.Email)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, userID uint, rawRefreshToken string) (string, string, error) {
	// Verify Refresh Token
	valid, err := s.authRepo.VerifyRefreshToken(ctx, userID, rawRefreshToken)
	if err != nil || !valid {
		return "", "", ErrInvalidToken
	}

	// Fetch New User Data
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", "", err
	}

	// Rotate Token
	return s.generateTokenPair(ctx, user.ID, user.Name, user.Username, user.Email)
}

func (s *AuthService) Logout(ctx context.Context, userID uint, jti string, accessTTL time.Duration) error {
	// Delete Refresh Token
	if err := s.authRepo.DeleteRefreshToken(ctx, userID); err != nil {
		return err
	}

	// Blacklist JTI Access Token
	if jti != "" && accessTTL > 0 {
		if err := s.authRepo.BlacklistAccessToken(ctx, jti, accessTTL); err != nil {
			return err
		}
	}

	return nil
}

// Access Token dan Refresh Token Pair Access Token dan Refresh Token Helper
func (s *AuthService) generateTokenPair(ctx context.Context, userID uint, name, username, email string) (string, string, error) {
	// --- Access Token ---
	accessJTI := uuid.New().String()
	accessClaims := map[string]any{
		"jti":      accessJTI,
		"user_id":  userID,
		"name":     name,
		"username": username,
		"email":    email,
	}
	jwtauth.SetExpiryIn(accessClaims, 15*time.Minute) // Access Token berumur pendek (15 Menit)
	jwtauth.SetIssuedNow(accessClaims)

	_, accessToken, err := s.accessTokenAuth.Encode(accessClaims)
	if err != nil {
		return "", "", err
	}

	// --- Refresh Token ---
	refreshTTL := 7 * 24 * time.Hour // 7 Hari
	refreshClaims := map[string]any{
		"jti":     uuid.New().String(),
		"user_id": userID,
	}
	jwtauth.SetExpiryIn(refreshClaims, refreshTTL)
	jwtauth.SetIssuedNow(refreshClaims)

	_, refreshToken, err := s.refreshTokenAuth.Encode(refreshClaims)
	if err != nil {
		return "", "", err
	}

	// Simpan Refresh Token di Storage (Redis + DB Fallback)
	if err := s.authRepo.StoreRefreshToken(ctx, userID, refreshToken, refreshTTL); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}