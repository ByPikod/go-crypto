package routes

import (
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/gofiber/fiber/v2"
)

type RegisterPaylaod struct {
	Name     string `json:"name"`
	Lastname string `json:"lastName"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

// Register account
// @Summary		Creates an account.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		name			body		string	true	"First name"	maxlength(32) minlength(2) example("John")
// @Param		lastName		body		string	true	"Last name"		maxlength(32) minlength(2) example("Doe")
// @Param		mail			body		string	true	"Mail address"	format(email) example("johndoe@example.com")
// @Param		password		body		string	true	"Password"		minlength(8) maxlength(256) example("JoHnDoe123")
// @Success		200				{object}	interface{}
// @Failure		400				{object}	interface{}
// @Router		/user/register	[post]
func Register(ctx *fiber.Ctx) error {

	// Parse payload
	payload := RegisterPaylaod{}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
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
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	// Bad request
	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	// Successfull
	return ctx.Status(200).JSON(result)

}

// @Summary		Login
// @Description Retrieves the token by credentials sent.
// @Tags		user
// @Accept		json
// @Produce		json
// @Param		mail			body		string	true	"Mail address"	format(email) example("johndoe@example.com")
// @Param		password		body		string	true	"Password"		minlength(8) maxlength(256) example("JoHnDoe123")
// @Success		200				{object}	interface{}
// @Failure		400				{object}	interface{}
// @Router		/user/login		[post]
func Login(ctx *fiber.Ctx) error {

	// Parse payload
	payload := new(models.UserLoginPayload)
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx, "Error: "+err.Error())
	}

	// Process
	result, err := payload.Login()

	if err != nil {
		helpers.LogError(err.Error())
		return helpers.InternalServerError(ctx)
	}

	// Bad request
	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	// Success
	return ctx.Status(200).JSON(result)

}

// @Summary		Login
// @Description Retrieves the user data.
// @Tags		user
// @Produce		json
// @Success		200				{object}	interface{}
// @Failure		401				{object}	interface{}
// @Security 	ApiKeyAuth
// @Router		/user/me		[get]
func Me(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*models.User)
	return ctx.Status(200).JSON(fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"lastname": user.Lastname,
		"mail":     user.Mail,
	})
}
