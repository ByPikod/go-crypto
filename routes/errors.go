package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Returns forbidden HTTP status code
func Forbidden(ctx *fiber.Ctx) error {
	return ctx.Status(403).JSON(fiber.Map{
		"code":    403,
		"message": "Forbidden",
	})
}

// Returns internal server error HTTP status code
func InternalServerError(ctx *fiber.Ctx) error {
	return ctx.Status(500).JSON(fiber.Map{
		"code":    500,
		"message": "Internal Server Error",
	})
}

// Returns bad request HTTP status code
func BadRequest(ctx *fiber.Ctx, reason ...string) error {
	message := "Bad Request"
	if len(reason) > 0 {
		message = strings.Join(reason, " ")
	}

	return ctx.Status(400).JSON(fiber.Map{
		"code":    400,
		"message": message,
	})
}

// Returns not found HTTP request code
func NotFound(ctx *fiber.Ctx) error {
	return ctx.Status(404).JSON(fiber.Map{
		"code":    404,
		"message": "Not found",
	})
}

// Returns unauthorized not found HTTP request code
func Unauthorized(ctx *fiber.Ctx, reason ...string) error {
	message := "Unauthorized"
	if len(reason) > 0 {
		message = strings.Join(reason, " ")
	}

	return ctx.Status(401).JSON(fiber.Map{
		"code":    401,
		"message": message,
	})
}
