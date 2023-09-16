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
	user := models.User{
		Name:     payload.Name,
		Lastname: payload.Lastname,
		Mail:     payload.Mail,
		Password: payload.Password,
	}

	result, err := user.Create()

	// An error occured
	if err != nil {
		return InternalServerError(ctx)
	}

	// Bad request
	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	// Successfull
	return ctx.Status(200).JSON(result)

}

// Handle "/api/user/login" endpoint
func Login(ctx *fiber.Ctx) error {

	// Parse payload
	payload := models.UserLoginCredentials{}
	if err := ctx.BodyParser(&payload); err != nil {
		return BadRequest(ctx, "Error: "+err.Error())
	}

	// Process
	result, err := payload.Login()

	if err != nil {
		return InternalServerError(ctx)
	}

	// Bad request
	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	// Success
	return ctx.Status(200).JSON(result)

}

func Me(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*models.User)
	return ctx.Status(200).JSON(fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"lastname": user.Lastname,
		"mail":     user.Mail,
	})
}
