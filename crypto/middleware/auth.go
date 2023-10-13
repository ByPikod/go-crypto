package middleware

import (
	"strings"

	"github.com/ByPikod/go-crypto/tree/crypto/helpers"
	"github.com/ByPikod/go-crypto/tree/crypto/log"
	"github.com/ByPikod/go-crypto/tree/crypto/services"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type (
	AuthMiddleware struct {
		userService *services.UserService
	}
)

func NewAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (authMiddleware *AuthMiddleware) Auth(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")

	// If token not passed in header
	if tokenString == "" {
		return helpers.Unauthorized(ctx)
	}

	// Token parse
	parts := strings.Split(tokenString, "Bearer ")
	if len(parts) > 1 {
		tokenString = parts[1]
	} else {
		return helpers.Unauthorized(ctx, "Token malformed.")
	}

	// Check token
	userID, err := authMiddleware.userService.Authenticate(tokenString)
	if err != nil {
		return helpers.Unauthorized(ctx)
	}

	// Look for user
	user, err := authMiddleware.userService.GetUserById(userID)
	if err != nil {
		// Database error
		return helpers.InternalServerError(ctx)
	}
	if user == nil {
		// User not found?
		log.Warn(
			"User sent a token that belongs to no account.",
			zap.String("remote_ip", ctx.IP()),
			zap.Uint("user_id", userID),
		)
		return helpers.Unauthorized(ctx, "User account removed or suspended!")
	}

	// Authenticated
	ctx.Locals("user", user)
	return ctx.Next()
}
