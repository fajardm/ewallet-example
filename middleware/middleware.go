package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fajardm/ewallet-example/errorcode"
	"github.com/gofiber/fiber"
	jwtware "github.com/gofiber/jwt"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"net/http"
)

func Protected() func(*fiber.Ctx) {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(viper.GetString("APP_SECRET")),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) {
	if err.Error() == "Missing or malformed JWT" {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT"})
	} else {
		c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT"})
	}
}

func GetUserID(ctx *fiber.Ctx) (*uuid.UUID, error) {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)
	id, err := uuid.FromString(userID)
	if err != nil {
		return nil, errorcode.ErrUnauthorized
	}
	return &id, nil
}
