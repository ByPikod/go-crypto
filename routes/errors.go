package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Returns forbidden HTTP status code
func Forbidden(c *fiber.Ctx) error {
	return c.Status(403).JSON(fiber.Map{
		"code":    403,
		"message": "Forbidden",
	})
}

// Returns internal server error HTTP status code
func InternalServerError(c *fiber.Ctx) error {
	return c.Status(500).JSON(fiber.Map{
		"code":    500,
		"message": "Internal Server Error",
	})
}

// Returns bad request HTTP status code
func BadRequest(c *fiber.Ctx, reason ...string) error {
	message := "Bad Request"
	if len(reason) > 0 {
		message = strings.Join(reason, " ")
	}

	return c.Status(500).JSON(fiber.Map{
		"code":    400,
		"message": message,
	})
}

// Returns not found HTTP request code
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).JSON(fiber.Map{
		"code":    404,
		"message": "Not found",
	})
}

// Returns unauthorized not found HTTP request code
func Unauthorized(c *fiber.Ctx) error {
	return c.Status(401).JSON(fiber.Map{
		"code":    401,
		"message": "Unauthorized",
	})
}
