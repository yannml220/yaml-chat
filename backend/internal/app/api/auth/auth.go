package auth

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenExpiration  = 30 * time.Minute
	RefreshTokenExpiration = 24 * time.Hour
	googleOauthConsentUrl  = "https://accounts.google.com/o/oauth2/v2/auth"
	googleOauthTokenUrl    = "https://oauth2.googleapis.com/token"
)

type Session struct {
	Id                    string     `json:"id"`
	UserId                string     `json:"user_id"`
	DeviceId              string     `json:"device_id"`
	UserAgent             string     `json:"user_agent"`
	HashedRefreshToken    string     `json:"hashed_refresh_token"`
	RefreshTokenExpiresAt *time.Time `json:"refresh_token_expires_at,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
	ExpiresAt             *time.Time `json:"expires_at,omitempty"`
	LastActivityAt        *time.Time `json:"last_activity_at,omitempty"`
}

type Identity struct {
	Id             string     `json:"id"`
	ProviderUserId string     `json:"provider_user_id"`
	UserId         string     `json:"user_id"`
	Provider       string     `json:"provider"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	LastSignInAt   *time.Time `json:"last_signin_at,omitempty"`
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	SessionId string `json:"session_id"`
	DeviceId  string `json:"device_id"`
	Role      string `json:"role"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	DeviceId  string `json:"device_id"`
	SessionId string `json:"session_id"`
}

type EmailTokenClaims struct {
	jwt.RegisteredClaims
}

type IdTokenClaims struct {
	Sub           string `json:"sub " validate:"required"`
	Email         string `json:"email" validate:"required"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name" validate:"required"`
	Exp           int64  `json:"exp" validate:"required"`
	Picture       string `json:"picture" validate:"required"`
}

type AuthContext struct {
	UserId    string `json:"user_id"`
	Role      string `json:"role"`
	SessionId string `json:"session_id"`
}
