package controllers

import (
	"errors"

	"github.com/ByPikod/go-crypto/tree/crypto/helpers"
	"github.com/ByPikod/go-crypto/tree/crypto/log"
	"github.com/ByPikod/go-crypto/tree/crypto/models"
	"github.com/ByPikod/go-crypto/tree/crypto/services"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type (
	UserController struct {
		service *services.UserService
	}
	RegisterPaylaod struct {
		Name         string `json:"name"`
		Lastname     string `json:"lastName"`
		Mail         string `json:"mail"`
		Password     string `json:"password"`
		Verification string `json:"verification"`
	}
	RegisterResponse struct {
		Status           bool   `json:"status"`
		Message          string `json:"message"`
		VerificationSent bool   `json:"verificationSent"`
	}
	UserLoginPayload struct {
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}
)

// Create user repository
func NewUserController(service *services.UserService) *UserController {
	return &UserController{service: service}
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
// @Success		200				{object}	RegisterResponse
// @Failure		400				{object}	interface{}
// @Router		/user/register	[post]
func (controller UserController) Register(ctx *fiber.Ctx) error {

	// Parse payload
	payload := RegisterPaylaod{}
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx)
	}

	// Validate data
	ok, reason := helpers.ValidateRegistration(
		payload.Name,
		payload.Lastname,
		payload.Mail,
		payload.Password,
	)

	if !ok {
		return ctx.Status(400).JSON(RegisterResponse{
			Status:           false,
			Message:          reason,
			VerificationSent: false,
		})
	}

	// Process
	err := controller.service.Create(
		payload.Name,
		payload.Lastname,
		payload.Mail,
		payload.Password,
		payload.Verification,
	)

	if err == nil {
		// No error, succesfully registered
		log.Info(
			"Account registered",
			zap.String("remote_ip", ctx.IP()),
			zap.String("name", payload.Name),
			zap.String("lastname", payload.Lastname),
			zap.String("user_mail", payload.Mail),
		)
		return ctx.Status(200).JSON(RegisterResponse{
			Status:           true,
			Message:          "OK",
			VerificationSent: false,
		})
	}

	if errors.Is(err, services.ErrUnavailableMailAddress) {
		// Mail is unavailable
		return ctx.Status(200).JSON(RegisterResponse{
			Status:           false,
			Message:          "Unavailable mail address.",
			VerificationSent: false,
		})
	}

	if errors.Is(err, services.ErrVerificationBlocked) {
		// Log
		log.Info(
			"E-mail verification blocked due to too many failed attempts. Try to receive another mail.",
			zap.String("remote_ip", ctx.IP()),
			zap.String("mail", payload.Mail),
		)

		return ctx.Status(200).JSON(RegisterResponse{
			Status:           false,
			Message:          "Email confirmation was blocked due to too many incorrect attempts!",
			VerificationSent: true,
		})
	}

	// Mail verification
	if errors.Is(err, services.ErrVerificationNeeded) {
		// Generate verification code.
		code, err := controller.service.CreateNewVerification(payload.Mail)

		if err != nil && errors.Is(err, services.ErrVerificationCooldown) {
			// Verification is still in cooldown
			return ctx.Status(200).JSON(RegisterResponse{
				Status:           false,
				Message:          "We just sent you a verification mail, you must wait for a while before requesting a new one.",
				VerificationSent: false,
			})
		}

		if err != nil {
			// Internal server error
			log.ControllerError("Register", err)
			return helpers.InternalServerError(ctx)
		}

		// connect kafka and put mail to queue
		// ...

		// Log
		log.Info(
			"Sent a mail verification code.",
			zap.String("remote_ip", ctx.IP()),
			zap.String("mail", payload.Mail),
			zap.String("code", code),
		)

		return ctx.Status(200).JSON(RegisterResponse{
			Status:           false,
			Message:          "We sent a code to your mail address!",
			VerificationSent: true,
		})
	}

	if errors.Is(err, services.ErrInvalidVerification) {
		// Verification code is incorrect
		return ctx.Status(200).JSON(RegisterResponse{
			Status:           false,
			Message:          "Verification code is incorrect!",
			VerificationSent: false,
		})
	}

	// Error is internal server error.
	log.ControllerError("Register", err)
	return helpers.InternalServerError(ctx)

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
func (controller *UserController) Login(ctx *fiber.Ctx) error {

	// Parse payload
	payload := new(UserLoginPayload)
	if err := ctx.BodyParser(&payload); err != nil {
		return helpers.BadRequest(ctx, "Error: "+err.Error())
	}

	// Process
	result, err := controller.service.Login(payload.Mail, payload.Password)

	if err != nil {
		log.ControllerError("Login", err)
		return helpers.InternalServerError(ctx)
	}

	// Bad request
	if result["status"] == false {
		return ctx.Status(400).JSON(result)
	}

	log.Info(
		"Client logged in",
		zap.String("remote_ip", ctx.IP()),
		zap.String("user_mail", payload.Mail),
	)

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
func (controller *UserController) Me(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*models.User)
	return ctx.Status(200).JSON(fiber.Map{
		"id":       user.ID,
		"name":     user.Name,
		"lastname": user.Lastname,
		"mail":     user.Mail,
	})
}
