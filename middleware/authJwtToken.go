package middleware

import (
	"os"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func VerificationAccessToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		TokenLookup:  "cookie:accessToken",
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: jwtError,
		ContextKey:   "access_token_context",
	})
}

func VerificationRefreshToken() fiber.Handler {
	return jwtware.New(jwtware.Config{
		TokenLookup:  "cookie:refreshToken",
		SigningKey:   jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		ErrorHandler: jwtError,
		ContextKey:   "refresh_token_context",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
