package routes

import (
	"net/mail"
	"time"

	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/gofiber/fiber/v2"
)

// Handle "/api/user/register" endpoint
func Register(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return BadRequest(ctx)
	}

	checkEmpty := user.Name == "" || user.Lastname == "" || user.Mail == "" || user.Password == ""
	if checkEmpty {
		return BadRequest(ctx)
	}

	// Validate username
	if len(user.Name) < 2 || len(user.Name) > 32 {
		return BadRequest(ctx, "The name must be between 2 and 32 characters in length.")
	}

	// Validate username
	if len(user.Lastname) < 2 || len(user.Lastname) > 32 {
		return BadRequest(ctx, "The lastname must be between 2 and 32 characters in length.")
	}

	// Validate password
	if isValid, reason := helpers.ValidatePassword(user.Password); !isValid {
		return BadRequest(ctx, reason)
	}

	// Validate mail address
	_, err := mail.ParseAddress(user.Mail)
	if err != nil {
		return BadRequest(ctx, "Mail address is not valid.")
	}

	exists, err := core.CheckExistsInDatabase(&models.User{Mail: user.Mail})
	if err != nil {
		return InternalServerError(ctx)
	}
	if exists {
		return BadRequest(ctx, "Mail address you specified is not available.")
	}

	// Create user
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err = core.DB.Create(&user).Error
	if err != nil {
		return InternalServerError(ctx)
	}

	ctx.Status(200).JSON(fiber.Map{
		"code":    "200",
		"message": "User registration successful",
	})

	return nil
}

// Handle "/api/user/login" endpoint
func Login(ctx *fiber.Ctx) error {

	ctx.Status(200).JSON(fiber.Map{
		"code":    "200",
		"message": "User registration successful",
	})

	return nil
}
