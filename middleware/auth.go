package middleware

import (
	"math"
	"strings"

	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")

	// If token not passed in header
	if tokenString == "" {
		return routes.Unauthorized(ctx)
	}

	// Token parse
	parts := strings.Split(tokenString, "Bearer ")
	if len(parts) > 1 {
		tokenString = parts[1]
	} else {
		return routes.Unauthorized(ctx, "Token malformed.")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(core.Config.AuthSecret), nil
	})

	if err != nil {
		// Failed to parse token
		return routes.Unauthorized(ctx, "Token malformed.")
	}

	// Get claims by decoding the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return routes.Unauthorized(ctx)
	}

	userID := claims["UserID"].(float64)
	if err != nil {
		return routes.Unauthorized(ctx, "Failed to parse ID.")
	}

	userID_uint := uint(math.Abs(float64(userID)))
	user, err := models.GetUserById(userID_uint)
	if err != nil {
		return routes.Unauthorized(ctx, "Token signature verified, but claimed user not found.")
	}
	if user == nil {
		return routes.Unauthorized(ctx, "User account removed or suspended!")
	}

	ctx.Locals("user", user)
	return ctx.Next()
}
