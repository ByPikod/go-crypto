package middleware

import "github.com/gofiber/fiber/v2"

func Json(ctx *fiber.Ctx) error {
	ctx.Accepts("application/json")
	return ctx.Next()
}
