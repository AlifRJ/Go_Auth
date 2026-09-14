package model

import (
	"context"
	"time"
)

type AuthSession struct {
	ID           uint      `json:"id" db:"id"`
	UserID       uint      `json:"user_id" db:"user_id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time 	`json:"created_at" db:"created_at"`
}

type TokenPayload struct {
	JTI       string    `json:"jti"`
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

type AuthRepository interface {
	// Storage Management untuk Refresh Token
	StoreRefreshToken(ctx context.Context, userID uint, token string, ttl time.Duration) error
	VerifyRefreshToken(ctx context.Context, userID uint, token string) (bool, error)
	DeleteRefreshToken(ctx context.Context, userID uint) error

	// Access Token Revocation / Blacklisting
	BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
}
