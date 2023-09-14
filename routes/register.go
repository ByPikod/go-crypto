package routes

import (
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return BadRequest(ctx)
	}

	checkEmpty := user.Name == "" || user.Lastname == "" || user.Mail == "" || user.Password == ""
	if checkEmpty {
		return BadRequest(ctx)
	}

	if len(user.Name) < 2 || len(user.Lastname) > 32 {
		return BadRequest(ctx, "The name must be between 2 and 32 characters in length.")
	}

	if isValid, reason := helpers.ValidatePassword(user.Password); !isValid {
		return BadRequest(ctx, reason)
	}

	ctx.Status(200).JSON(fiber.Map{
		"code":    "200",
		"message": "User registration successful",
	})

	return nil
}
