package models

import (
	"errors"
	"net/mail"

	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type User struct {
	BaseModel
	Name     string `json:"name" gorm:"not null"`
	Lastname string `json:"lastName" gorm:"not null"`
	Mail     string `json:"mail" gorm:"index;not null"`
	Password string `json:"password" gorm:"not null"`
}

// Registers a user.
func UserSignUp(
	name string,
	lastName string,
	mailAddress string,
	password string,
) (map[string]interface{}, error) {
	var checkEmpty bool = (name == "" || lastName == "" || mailAddress == "" || password == "")
	if checkEmpty {
		return fiber.Map{
			"status":  false,
			"message": "Blank fields.",
		}, nil
	}

	// Validate username
	if len(name) < 2 || len(name) > 32 {
		return fiber.Map{
			"status":  false,
			"message": "The name must be between 2 and 32 characters in length.",
		}, nil
	}

	// Validate username
	if len(lastName) < 2 || len(lastName) > 32 {
		return fiber.Map{
			"status":  false,
			"message": "The lastname must be between 2 and 32 characters in length.",
		}, nil
	}

	// Validate password
	if isValid, reason := helpers.ValidatePassword(password); !isValid {
		return fiber.Map{
			"status":  false,
			"message": reason,
		}, nil
	}

	// Validate mail address
	_, err := mail.ParseAddress(mailAddress)
	if err != nil {
		return fiber.Map{
			"status":  false,
			"message": "Mail address is not valid.",
		}, nil
	}

	// Check if mail address exists in db
	exists, err := core.CheckExistsInDatabase(&User{Mail: mailAddress})
	if err != nil {
		return nil, err
	}
	if exists {
		return fiber.Map{
			"status":  false,
			"message": "Mail address you specified is not available.",
		}, nil
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Name:     name,
		Lastname: lastName,
		Mail:     mailAddress,
		Password: hashedPassword,
	}

	err = core.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"status":  true,
		"message": "User registration completed.",
	}, nil
}

// User login
func UserSignIn(mailAddress string, password string) (map[string]interface{}, error) {

	const incorrectCredentials = "Incorrect mail adress or password."
	if mailAddress == "" || password == "" {
		return fiber.Map{
			"status":  false,
			"message": incorrectCredentials,
		}, nil
	}

	/*
		// Validate mail address
		_, err := mail.ParseAddress(mailAddress)
		if err != nil {
			return BadRequest(ctx, incorrectCredentials)
		}

		// Validate password
		if isValid, _ := helpers.ValidatePassword(password); !isValid {
			return BadRequest(ctx, incorrectCredentials)
		}
	*/

	// Find the matching mail address from the database.
	matching := User{}
	res := core.DB.Model(&matching).Where(&User{Mail: mailAddress}).First(&matching)
	if res.Error != nil {
		// If matching user credentials not found
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return fiber.Map{
				"status":  false,
				"message": incorrectCredentials,
			}, nil
		}
		// If something else happened
		return nil, res.Error
	}

	// Check if matching user password and payload password equals
	passed, err := helpers.ComparePasswords(matching.Password, password)
	if err != nil {
		return nil, err
	}
	if !passed {
		return fiber.Map{
			"status":  false,
			"message": incorrectCredentials,
		}, nil
	}

	return fiber.Map{
		"status":  true,
		"message": "OK!",
	}, nil

}
