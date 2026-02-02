package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/yannml220/chat_agent_app/internal/app/api/user"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
)

type AuthService interface {
	SigninWithOauth_(ctx context.Context, deviceId string, user_agent string, idToken *IdTokenClaims, provider string) (string, string, error)

	IsTokenReuseDetected(ctx context.Context, hashedCookieRefreshToken string, sessionId string, userId string, deviceId string, userAgent string) (bool, error)

	ValidateTokenPairWithReuseDetection(ctx context.Context, refreshToken string, accessToken string, deviceId string, userAgent string) (*TokenValidationResult, error)

	RotateRefreshTokenAndExtendSession(ctx context.Context, oldHashedToken string, userId string, sessionId string, deviceId string, userAgent string, newHashedToken string, newExpiresAt time.Time) error

	GetUserIdentity(ctx context.Context, userId string, provider string) (*Identity, error)

	CreateIdentity(ctx context.Context, identity *Identity, userId string) (string, error)

	GetSessionById(ctx context.Context, id string) (*Session, error)

	GetSessionByRefreshToken(ctx context.Context, hashedRefreshToken string) (*Session, error)

	CreateSession_(ctx context.Context, session *Session, userId string) (string, error)

	DeleteSessionById(ctx context.Context, sessionId string) error

	DeleteSessionByRefreshTokenHash(ctx context.Context, hashed_refresh_token string) error

	InvalidateAllUserSessions(ctx context.Context, userId string) error

	HashPassword(ctx context.Context, password string, cost int) (string, error)
	HashToken(ctx context.Context, token string, key string) string
	VerifyTokenHash(ctx context.Context, hashedToken string, token string, key string) bool
	VerifyPassword(ctx context.Context, hashedPassword string, password string) error

	DecodeIdToken(ctx context.Context, idToken string, clientID string) (*IdTokenClaims, error)
	VerifyRefreshToken(ctx context.Context, tokenString string) (*jwt.Token, error)
	VerifyAccessToken(ctx context.Context, tokenString string) (*jwt.Token, error)
	VerifyEmailToken(ctx context.Context, tokenString string) (*jwt.Token, error)

	GenerateToken(ctx context.Context, signingMethod jwt.SigningMethod, claims jwt.Claims, key string) (string, error)

	GetUserByEmail(ctx context.Context, email string) (*user.User, error)

	UpdateUserIdentity(ctx context.Context, identity *Identity, userId string, provider string) error

	VerifyUserEmail(ctx context.Context, userId string) error

	CreateUserAndIdentity(ctx context.Context, user *user.User, identity *Identity) (string, string, error)
}

type AuthServiceImpl struct {
	types.Service
	userService user.UserService
	repo        AuthRepo
}

func NewService(repo AuthRepo, userService user.UserService) AuthService {

	return &AuthServiceImpl{
		types.Service{},
		userService,
		repo,
	}

}

func (a *AuthServiceImpl) RotateRefreshTokenAndExtendSession(ctx context.Context, oldHashedToken string, userId string, sessionId string, deviceId string, userAgent string, newHashedToken string, newExpiresAt time.Time) error {
	return a.repo.RotateRefreshTokenAndExtendSession(ctx, oldHashedToken, userId, sessionId, deviceId, userAgent, newHashedToken, newExpiresAt)

}

func (a *AuthServiceImpl) CreateSession_(ctx context.Context, session *Session, userId string) (string, error) {
	return a.repo.CreateSession_(ctx, session, userId)
}

type TokenValidationResult struct {
	UserId                string
	SessionId             string
	OldHashedRefreshToken string
	IsReuse               bool
}

func (a *AuthServiceImpl) ValidateTokenPairWithReuseDetection(ctx context.Context, refreshToken string, accessToken string, deviceId string, userAgent string) (*TokenValidationResult, error) {

	if deviceId == "" || userAgent == "" {
		return &TokenValidationResult{}, fmt.Errorf("deviceId or userAgent missing")
	}

	log.Print("access token : ", accessToken)
	log.Print("refresh token : ", refreshToken)
	if refreshToken == "" {
		return &TokenValidationResult{}, fmt.Errorf("missing refresh token")
	}

	if accessToken == "" {
		return &TokenValidationResult{}, fmt.Errorf("missing access token")
	}

	rtoken, err := a.VerifyRefreshToken(ctx, refreshToken)
	if err != nil {
		return &TokenValidationResult{}, err
	}

	atoken, err := a.VerifyAccessToken(ctx, accessToken)
	if err != nil {
		return &TokenValidationResult{}, err
	}

	rtSub, _ := rtoken.Claims.GetSubject()
	atSub, _ := atoken.Claims.GetSubject()

	hashedToken := a.HashToken(ctx, refreshToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))

	if rtSub != atSub {
		return &TokenValidationResult{
			UserId:                rtSub,
			OldHashedRefreshToken: hashedToken,
			IsReuse:               true,
		}, nil

	}

	claims, ok := atoken.Claims.(*AccessTokenClaims)
	if !ok || !atoken.Valid {
		return &TokenValidationResult{}, fmt.Errorf("invalid claims")
	}

	rtClaims, ok := rtoken.Claims.(*RefreshTokenClaims)
	if !ok || !rtoken.Valid {
		return &TokenValidationResult{}, fmt.Errorf("invalid refresh claims")
	}

	if rtClaims.DeviceId != deviceId {
		return &TokenValidationResult{
			UserId:                rtSub,
			OldHashedRefreshToken: hashedToken,
			IsReuse:               true,
		}, nil
	}

	isReuse, err := a.IsTokenReuseDetected(ctx, hashedToken, claims.SessionId, rtSub, deviceId, userAgent)
	if err != nil {
		return &TokenValidationResult{}, err
	}

	return &TokenValidationResult{
		UserId:                rtSub,
		SessionId:             claims.SessionId,
		OldHashedRefreshToken: hashedToken,
		IsReuse:               isReuse,
	}, nil
}

func (a *AuthServiceImpl) IsTokenReuseDetected(ctx context.Context, hashedCookieRefreshToken string, sessionId string, userId string, deviceId string, userAgent string) (bool, error) {
	return a.repo.IsTokenReuseDetected(ctx, hashedCookieRefreshToken, sessionId, userId, deviceId, userAgent)

}

func (a *AuthServiceImpl) GetSessionByRefreshToken(ctx context.Context, hashedRefreshToken string) (*Session, error) {
	return a.repo.GetSessionByRefreshToken(ctx, hashedRefreshToken)
}

func (a *AuthServiceImpl) GetSessionById(ctx context.Context, id string) (*Session, error) {
	return a.repo.GetSessionById(ctx, id)
}

func (a *AuthServiceImpl) DeleteSessionByRefreshTokenHash(ctx context.Context, hashed_refresh_token string) error {
	return a.repo.DeleteSessionByRefreshTokenHash(ctx, hashed_refresh_token)
}

func (a *AuthServiceImpl) DeleteSessionById(ctx context.Context, sessionId string) error {
	return a.repo.DeleteSessionById(ctx, sessionId)
}

func (a *AuthServiceImpl) HashPassword(ctx context.Context, password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err

	}
	return string(hashedPassword), nil

}

func (a *AuthServiceImpl) VerifyPassword(ctx context.Context, hashedPassword string, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}

func (a *AuthServiceImpl) HashToken(ctx context.Context, token string, key string) string {

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))

}

func (a *AuthServiceImpl) VerifyTokenHash(ctx context.Context, hashedToken string, token string, key string) bool {
	return hashedToken == a.HashToken(ctx, token, key)
}

func (a *AuthServiceImpl) InvalidateAllUserSessions(ctx context.Context, userId string) error {

	return a.repo.DeleteAllUserSessions(ctx, userId)

}

func (a *AuthServiceImpl) GenerateToken(ctx context.Context, signingMethod jwt.SigningMethod, claims jwt.Claims, key string) (string, error) {

	token := jwt.NewWithClaims(signingMethod, claims)
	var tokenString string
	var err error

	if tokenString, err = token.SignedString([]byte(key)); err != nil {
		return "", err
	}

	return tokenString, nil

}

func (a *AuthServiceImpl) VerifyRefreshToken(
	ctx context.Context,
	tokenString string,
) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
			}

			key := os.Getenv("REFRESH_TOKEN_SECRET")
			if key == "" {
				return nil, errors.New("clé secrète manquante")
			}
			return []byte(key), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("refresh token invalide")
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		return nil, errors.New("claims du refresh token invalides")
	}

	if claims.DeviceId == "" {
		return nil, errors.New("device_id manquant")
	}

	return token, nil
}

func (a *AuthServiceImpl) VerifyAccessToken(
	ctx context.Context,
	tokenString string,
) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
			}

			key := os.Getenv("ACCESS_TOKEN_SECRET")
			if key == "" {
				return nil, errors.New("clé secrète manquante dans l'environnement")
			}
			return []byte(key), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return token, nil
		}

		return nil, err
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide")
	}

	if claims.SessionId == "" {
		return nil, errors.New("session_id manquant")
	}

	return token, nil
}

func (a *AuthServiceImpl) VerifyEmailToken(ctx context.Context, tokenString string) (*jwt.Token, error) {

	token, err := jwt.ParseWithClaims(tokenString, &EmailTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// verify the signature method

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("It's not the expected  signature method")
		}
		// get the key from env

		key := os.Getenv("EMAIL_TOKEN_SECRET")

		if key == "" {

			return nil, errors.New("Couldn't load the key")

		}

		return []byte(key), nil

	})

	// we couldn't parse
	if err != nil {
		log.Print("we stopped here")
		return nil, err
	}

	// we parsed but it's invalid
	if !token.Valid {
		return nil, errors.New("Invalid token")

	}

	return token, nil

	//verify the claims validity

}

func (a *AuthServiceImpl) DecodeIdToken(ctx context.Context, idToken string, clientID string) (*IdTokenClaims, error) {

	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return nil, err
	}

	// Verify the token
	payload, err := validator.Validate(ctx, idToken, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	claims := IdTokenClaims{
		Sub:     payload.Subject,
		Email:   payload.Claims["email"].(string),
		Name:    payload.Claims["name"].(string),
		Exp:     payload.Expires,
		Picture: payload.Claims["picture"].(string),
	}

	return &claims, nil

}

func (a *AuthServiceImpl) CreateIdentity(ctx context.Context, identity *Identity, userId string) (string, error) {

	return a.repo.CreateIdentity(ctx, identity, userId)

}

func (a *AuthServiceImpl) GetUserIdentity(ctx context.Context, userId string, provider string) (*Identity, error) {
	return a.repo.GetUserIdentity(ctx, userId, provider)

}

func (a *AuthServiceImpl) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {

	return a.repo.GetUserByEmail(ctx, email)

}

func (a *AuthServiceImpl) UpdateUserIdentity(ctx context.Context, identity *Identity, userId string, provider string) error {

	return a.repo.UpdateUserIdentity(ctx, identity, userId, provider)

}

func (a *AuthServiceImpl) VerifyUserEmail(ctx context.Context, userId string) error {
	return a.repo.VerifyUserEmail(ctx, userId)

}

func (a *AuthServiceImpl) CreateUserAndIdentity(ctx context.Context, user *user.User, identity *Identity) (string, string, error) {
	return a.repo.CreateUserAndIdentity(ctx, user, identity)
}

func (a *AuthServiceImpl) SigninWithOauth_(ctx context.Context, deviceId string, userAgent string, idToken *IdTokenClaims, provider string) (string, string, error) {

	var userId string
	emptyIdentity := Identity{}
	var err error

	log.Print("id token sub :", idToken.Sub)
	log.Print("id token email :", idToken.Email)
	log.Print("id token confirmed :", idToken.EmailVerified)
	log.Print("id token profile picture :", idToken.Picture)
	log.Print("id token name :", idToken.Name)

	retrievedUser, err := a.userService.GetUserByEmail(ctx, idToken.Email)
	if err != nil {
		if err.Error() != "no rows in result set"  {
			log.Print("error retrieving user given email ", err.Error())
			return "", "", err

		}
	}

	if retrievedUser == nil {

		log.Print("iam entering the case 1 ( create user + identity )")
		newUser := &user.User{
			Email:             &idToken.Email,
			Name:              &idToken.Name,
			ProfilePictureUrl: &idToken.Picture,
			Confirmed:         &idToken.EmailVerified,
		}

		newIdentity := &Identity{
			ProviderUserId: idToken.Sub,
			Provider:       provider,
		}

		userId, _, err = a.CreateUserAndIdentity(ctx, newUser, newIdentity)
		if err != nil {
			return "", "", fmt.Errorf("couldn't create user and identity (case 1): %w", err)
		}
	} else {
		// Existing user
		userId = retrievedUser.Id

		updatedUser := &user.User{
			Email:             &idToken.Email,
			Confirmed:         &idToken.EmailVerified,
			Name:              &idToken.Name,
			ProfilePictureUrl: &idToken.Picture,
		}

		err = a.userService.UpdateUser(ctx, updatedUser, userId)

		if err != nil {
			return "", "", fmt.Errorf("couldn't update user (case 2): %w", err)
		}

		identity, err := a.repo.GetUserIdentity(ctx, userId, provider)
		if err != nil {
			return "", "", fmt.Errorf("couldn't get user's identity (case 2): %w", err)
		}

		if *identity == emptyIdentity {

			log.Print("iam entering the case 2 (Existing user, first time connecting via Google) ")

			newIdentity := &Identity{
				ProviderUserId: idToken.Sub,
				UserId:         userId,
				Provider:       provider,
			}

			_, err = a.CreateIdentity(ctx, newIdentity, userId)
			if err != nil {
				return "", "", fmt.Errorf("couldn't create identity (case 2a): %w", err)
			}
		} else {

			log.Print("iam entering the case 3 (Returning existing google user")

			updatedIdentity := &Identity{}

			err = a.UpdateUserIdentity(ctx, updatedIdentity, userId, provider)
			if err != nil {
				return "", "", fmt.Errorf("couldn't update identity (case 2b): %w", err)
			}
		}
	}

	// All cases: Create session + refresh token (Transaction B)
	newExpiration := time.Now().Add(RefreshTokenExpiration)

	rtClaims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(newExpiration),
		},
		DeviceId: deviceId,
	}

	rToken, err := a.GenerateToken(ctx, jwt.SigningMethodHS256, rtClaims, os.Getenv("REFRESH_TOKEN_SECRET"))

	if err != nil {
		return "", "", fmt.Errorf("couldn't generate refresh token: %w", err)
	}

	hashedToken := a.HashToken(ctx, rToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))

	log.Print(" Here is the hashed refresh token : ", hashedToken)

	newSession := &Session{
		DeviceId:              deviceId,
		UserAgent:             userAgent,
		ExpiresAt:             &newExpiration,
		HashedRefreshToken:    hashedToken,
		RefreshTokenExpiresAt: &newExpiration,
	}

	sessionId, err := a.CreateSession_(ctx, newSession, userId)

	if err != nil {
		return "", "", fmt.Errorf("couldn't create the session : %w", err)
	}

	// Generate access token
	atClaims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
		},
		SessionId: sessionId,
		DeviceId:  deviceId,
		Role:      "user",
	}

	aToken, err := a.GenerateToken(ctx, jwt.SigningMethodHS256, atClaims, os.Getenv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return "", "", fmt.Errorf("couldn't generate access token: %w", err)
	}

	return aToken, rToken, nil
}
