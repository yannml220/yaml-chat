package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"github.com/yannml220/chat_agent_app/internal/pkg/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthServerTokenResponse struct {
	AccessToken  string `json:"access_token" validate:"required"`
	ExpiresIn    int32  `json:"expires_in" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
	Scope        string `json:"scope" validate:"required"`
	TokenType    string `json:"token_type" validate:"required"`
	IdToken      string `json:"id_token" validate:"required"`
}

type AuthHandler interface {
	SignoutHandler_(c *fiber.Ctx) error
	RotateTokenHandler(c *fiber.Ctx) error
	OauthSigninHandler(c *fiber.Ctx) error
	Oauth2CallbackHandler(c *fiber.Ctx) error
	GetUserInfosHandler(c *fiber.Ctx) error
}

type authHandler struct {
	types.Handler
	Service AuthService
}

func NewHandler(as AuthService) AuthHandler {
	return &authHandler{
		Service: as,
	}

}

func (ah *authHandler) SignoutHandler_(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	refreshToken := c.Cookies("refresh_token")
	accessToken := c.Cookies("access_token")
	userAgent := c.Get("User-Agent")
	deviceId := c.Get("X-Device-ID")

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "id_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
	})

	c.ClearCookie("refresh_token")
	c.ClearCookie("access_token")
	c.ClearCookie("id_token")

	res, err := ah.Service.ValidateTokenPairWithReuseDetection(ctx, refreshToken, accessToken, deviceId, userAgent)

	if res.IsReuse {
		ah.Service.InvalidateAllUserSessions(ctx, res.UserId)
		return ah.JSONResponse(c, fiber.StatusNoContent, fiber.Map{})
	}
	if err != nil {
		if refreshToken != "" {

			hashedToken := ah.Service.HashToken(ctx, refreshToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))
			ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)
			return ah.JSONResponse(c, fiber.StatusNoContent, fiber.Map{})

		}

	}

	return ah.JSONResponse(c, fiber.StatusNoContent, fiber.Map{})
}

func (ah *authHandler) RotateTokenHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	refreshToken := c.Cookies("refresh_token")
	accessToken := c.Cookies("access_token")

	deviceId := c.Get("X-Device-ID")
	userAgent := c.Get("User-Agent")

	c.ClearCookie("refresh_token")
	c.ClearCookie("access_token")

	hashedToken := ah.Service.HashToken(ctx, refreshToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))

	res, err := ah.Service.ValidateTokenPairWithReuseDetection(ctx, refreshToken, accessToken, deviceId, userAgent)

	if err != nil && err.Error() != "missing access token" {
		log.Print("error : ", err.Error())
		if refreshToken != "" {
			ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)
		}
		return c.Redirect("/auth/signin", 401)
	}

	if res.IsReuse {
		ah.Service.InvalidateAllUserSessions(ctx, res.UserId)
		return c.Redirect("/auth/signin", 403)
	}

	session, err := ah.Service.GetSessionByRefreshToken(ctx, hashedToken)

	if err != nil || session == nil {

		ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)
		return c.Redirect("/auth/signin", 500)
	}

	rtClaims := RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   session.UserId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration)),
		},
		DeviceId: deviceId,
	}

	atClaims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   session.UserId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
		},
		SessionId: session.Id,
		DeviceId:  deviceId,
		Role:      "user",
	}

	// Generate a new refresh token
	newRToken, err := ah.Service.GenerateToken(ctx, jwt.SigningMethodHS256, rtClaims, os.Getenv("REFRESH_TOKEN_SECRET"))

	if err != nil {
		ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)

		return c.Redirect("/auth/signin", 500)
	}

	newAToken, err := ah.Service.GenerateToken(ctx, jwt.SigningMethodHS256, atClaims, os.Getenv("ACCESS_TOKEN_SECRET"))

	if err != nil {
		log.Print("error : ", err.Error())
		ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)

		return c.Redirect("/auth/signin", 500)
	}

	newHashedToken := ah.Service.HashToken(ctx, newRToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))

	// rotate the token and extend the associated session
	err = ah.Service.RotateRefreshTokenAndExtendSession(ctx, hashedToken, session.UserId, session.Id, deviceId, userAgent, newHashedToken, time.Now().Add(RefreshTokenExpiration))

	if err != nil {
		log.Print("error rotating the refresh token : ", err.Error())
		ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)

		return c.Redirect("/auth/signin", 500)
	}

	rtCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    newRToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(RefreshTokenExpiration),
	}

	atCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    newAToken,
		Path:     "/",
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(AccessTokenExpiration),
	}

	c.Cookie(rtCookie)
	c.Cookie(atCookie)

	return ah.JSONResponse(c, fiber.StatusOK, fiber.Map{})

}

// GET
func (ah *authHandler) Oauth2CallbackHandler(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 15*time.Second)

	defer cancel()

	refreshToken := c.Cookies("refresh_token")
	accessToken := c.Cookies("access_token")

	c.ClearCookie("refresh_token")
	c.ClearCookie("access_token")

	userAgent := c.Get("User-Agent")

	error := c.Query("error")
	if error != "" {

		log.Print("error : ", error)
		return ah.JSONResponse(c, fiber.ErrUnauthorized.Code, fiber.Map{
			"message":      "there is an error in the query param",
			"redirect_url": "http://localhost:3000/login",
		})
	}

	code := c.Query("code")

	if code == "" {

		return ah.JSONResponse(c, fiber.ErrUnauthorized.Code, fiber.Map{
			"message":      "the auth code is empty",
			"redirect_url": "http://localhost:3000/login",
		})
	}

	state := c.Query("state")
	if state == "" {
		return c.Status(400).SendString("State manquant")
	}

	stateParts := strings.Split(state, "|")
	if len(stateParts) != 2 {
		return ah.JSONResponse(c, fiber.ErrUnauthorized.Code, fiber.Map{
			"message":      "invalid state format",
			"redirect_url": "http://localhost:3000/login",
		})
	}
	receivedNonce := stateParts[0]
	deviceId := stateParts[1]

	cookieNonce := c.Cookies("oauth_state")
	if cookieNonce == ""  {
		log.Print("missing oauth_state cookie")
		return ah.JSONResponse(c, fiber.ErrUnauthorized.Code, fiber.Map{
			"message":      "missing oauth state",
			"redirect_url": "http://localhost:3000/login",
		})
	}

	if  cookieNonce != receivedNonce {
		log.Print("Invalid nonce")
		return ah.JSONResponse(c, fiber.ErrUnauthorized.Code, fiber.Map{
			"message":      "invalid nonce",
			"redirect_url": "http://localhost:3000/login",
		})
	}


	c.ClearCookie("oauth_state")

	if refreshToken != "" || accessToken != "" {

		res, err := ah.Service.ValidateTokenPairWithReuseDetection(ctx, refreshToken, accessToken, deviceId, userAgent)

		if err != nil {
			if refreshToken != "" {

				log.Print("error : ", err.Error())
				hashedToken := ah.Service.HashToken(ctx, refreshToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))
				ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)

			}
		}

		if res.IsReuse {
			ah.Service.InvalidateAllUserSessions(ctx, res.UserId)
		}
	}

	values := url.Values{}
	values.Set("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	values.Set("client_secret", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))

	resp, err := http.PostForm(googleOauthTokenUrl, values)

	if err != nil {

		log.Print("error : ", err.Error())
		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error":        err.Error(),
			"message":      "can't make the post request google token url",
			"redirect_url": "http://localhost:5173/login",
		})
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error":        fmt.Sprintf("Unexpected status code: %d", resp.StatusCode),
			"message":      "failed to get token from Google oauth server",
			"redirect_url": "http://localhost:5173/login",
		})
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Print("error : ", err.Error())
		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error":        err.Error(),
			"message":      "can't read the response body",
			"redirect_url": "http://localhost:5173/login",
		})

	}

	respStruct := new(AuthServerTokenResponse)

	err = json.Unmarshal(respBody, &respStruct)

	if err != nil {
		log.Print("error : ", err.Error())
		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error":        err.Error(),
			"message":      "can't unmarshal the google token response body",
			"redirect_url": "http://localhost:5173/login",
		})

	}

	// decode the id token
	idTokenClaims, err := ah.Service.DecodeIdToken(ctx, respStruct.IdToken, os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	if err != nil {

		log.Print("error : ", err.Error())
		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"redirect_url": "http://localhost:5173/login",
			"error":        err.Error(),
			"message":      "can't decode the id token",
		})

	}

	// appeler le service signin with oauth qui retourne les deux tokens

	aToken, rToken, err := ah.Service.SigninWithOauth_(ctx, deviceId, userAgent, idTokenClaims, "google")

	if err != nil {
		log.Print("error : ", err.Error())
		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error":        err.Error(),
			"message":      "can't sign in",
			"redirect_url": "http://localhost:5173/login",
		})
	}

	refreshTokenCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    rToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(RefreshTokenExpiration),
	}

	accessTokenCookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    aToken,
		Path:     "/",
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
		Expires:  time.Now().Add(AccessTokenExpiration),
	}

	idtKind := &fiber.Cookie{
		Name:     "idt_kind",
		Value:    "oauth",
		Path:     "/",
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
	}

	c.Cookie(refreshTokenCookie)
	c.Cookie(accessTokenCookie)
	c.Cookie(idtKind)

	return c.Redirect("http://localhost:5173/")

}

func (ah *authHandler) GetUserInfosHandler(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	authContext := c.Locals("auth_context")

	a, ok := authContext.(AuthContext)
	if !ok {
		log.Print("problem with the type of the auth context")
		return c.Status(500).SendString("auth context type error")
	}

	userId := a.UserId

	user, err := ah.Service.(*AuthServiceImpl).userService.GetUserById(ctx, userId)

	if err != nil {

		log.Print("(case 1) failed to fetch user by id :", err)

		return c.Redirect("http://localhost:3000/login", 401)
	}

	return ah.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"id":     userId,
		"name":   user.Name,
		"avatar": user.ProfilePictureUrl,
	})

}

func (ah *authHandler) OauthSigninHandler(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

	defer cancel()

	deviceId := c.Get("X-Device-ID")
	userAgent := c.Get("User-Agent")

	refreshToken := c.Cookies("refresh_token")
	accessToken := c.Cookies("access_token")

	c.ClearCookie("refresh_token")
	c.ClearCookie("access_token")

	if refreshToken != "" || accessToken != "" {

		hashedToken := ah.Service.HashToken(ctx, refreshToken, os.Getenv("REFRESH_TOKEN_HASH_SECRET"))

		res, err := ah.Service.ValidateTokenPairWithReuseDetection(ctx, refreshToken, accessToken, deviceId, userAgent)

		if err != nil {
			if refreshToken != "" {
				ah.Service.DeleteSessionByRefreshTokenHash(ctx, hashedToken)
			}
		}

		if res.IsReuse {
			ah.Service.InvalidateAllUserSessions(ctx, res.UserId)
		}
	}

	baseUrl := "https://accounts.google.com/o/oauth2/v2/auth"

	options := url.Values{}
	options.Add("redirect_uri", os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))
	options.Add("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	options.Add("access_type", "offline")
	options.Add("response_type", "code")

	nonce := make([]byte, 32)
	if _, err := rand.Read(nonce); err != nil {
		return ah.JSONResponse(c, fiber.ErrInternalServerError.Code, fiber.Map{
			"error": "failed to generate nonce",
		})
	}
	nonceStr := base64.RawURLEncoding.EncodeToString(nonce)

	state := fmt.Sprintf("%s|%s", nonceStr, deviceId)

	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    nonceStr,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	options.Add("state", state)
	options.Add("prompt", "consent")

	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	options.Add("scope", strings.Join(scopes, " "))

	url := fmt.Sprintf("%s?%s", baseUrl, options.Encode())

	return ah.JSONResponse(c, fiber.StatusOK, fiber.Map{
		"url": url,
	})

}
