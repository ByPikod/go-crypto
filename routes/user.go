package routes

import (
	"github.com/ByPikod/go-crypto/models"
	"github.com/gofiber/fiber/v2"
)

// Handle "/api/user/register" endpoint
func Register(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Name     string `json:"name"`
		Lastname string `json:"lastName"`
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx)
	}

	// Process
	result, err := models.UserSignUp(
		payload.Name,
		payload.Lastname,
		payload.Mail,
		payload.Password,
	)

	if err != nil {
		return InternalServerError(ctx)
	}

	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	ctx.Status(200).JSON(result)

	return nil
}

// Handle "/api/user/login" endpoint
func Login(ctx *fiber.Ctx) error {

	// Parse payload
	var payload struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx, "Error: "+err.Error())
	}

	// Process
	result, err := models.UserSignIn(
		payload.Mail,
		payload.Password,
	)

	if err != nil {
		return InternalServerError(ctx)
	}

	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	ctx.Status(200).JSON(result)

	return nil
}
