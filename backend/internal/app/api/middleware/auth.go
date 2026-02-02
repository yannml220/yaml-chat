package middleware

import (
	"context"
	"log"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/yannml220/chat_agent_app/internal/app/api/auth"
)

type AuthMiddleware struct {
	AuthService auth.AuthService
}

func NewAuthMiddleware(authService auth.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}

}

func (am *AuthMiddleware) Run() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)

		defer cancel()

		log.Print("Protected route !")

		accessToken := c.Cookies("access_token")
		log.Print("ACCESS TOK :", accessToken)
		if accessToken == "" {
			log.Print("Missing access token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing access token"})
		}

		deviceId := c.Get("X-Device-ID")
		if deviceId == "" {
			log.Print("X-Device-ID header required")
			c.ClearCookie("access_token")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "X-Device-ID header required"})
		}

		atoken, err := am.AuthService.VerifyAccessToken(ctx, accessToken)

		if err != nil {
			log.Print("Invalid access token")
			c.ClearCookie("access_token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid access token"})
		}

		claims, ok := atoken.Claims.(*auth.AccessTokenClaims)
		if !ok {
			log.Print("Invalid claims")
			c.ClearCookie("access_token")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid claims"})

		}

		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			log.Print("Access token expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Access token expired",
			})
		}

		if claims.DeviceId != deviceId {
			log.Print("Device mismatch")
			c.ClearCookie("access_token")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Device mismatch"})
		}

		authContext := auth.AuthContext{
			UserId:    claims.Subject,
			Role:      "user",
			SessionId: claims.SessionId,
		}

		c.Locals("auth_context", authContext)

		log.Print("Middleware checks passed !")

		return c.Next()
	}
}
