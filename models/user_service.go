package models

import (
	"errors"
	"net/mail"

	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

/*
 * Structs
 */

type UserLoginPayload struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

/*
 * Functions
 */

// Retrieves user by its ID, returns nil if user not found.
func GetUserById(id uint) (*User, error) {
	ret := User{}
	ret.ID = id
	result := core.DB.Model(&ret).Where(&ret).First(&ret)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &ret, nil
}

// User login
func (payload *UserLoginPayload) Login() (map[string]interface{}, error) {

	const incorrectCredentials = "Incorrect credentials!"
	// Check blank fields
	if payload.Mail == "" || payload.Password == "" {
		return fiber.Map{
			"status":  false,
			"message": "Blank fields.",
		}, nil
	}

	// Find the matching mail address from the database.
	matching := User{}
	res := core.DB.Model(&matching).Where(&User{Mail: payload.Mail}).First(&matching)
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
	passed, err := helpers.ComparePasswords(matching.Password, payload.Password)
	if err != nil {
		return nil, err
	}
	if !passed {
		return fiber.Map{
			"status":  false,
			"message": incorrectCredentials,
		}, nil
	}

	// Create token
	claims := (struct {
		jwt.StandardClaims
		UserID uint
	}{UserID: matching.ID})
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &claims)
	tokenString, err := token.SignedString([]byte(core.Config.AuthSecret))
	if err != nil {
		return nil, err
	}

	// Return token
	return fiber.Map{
		"status":  true,
		"message": "OK!",
		"token":   tokenString,
	}, nil

}

// Registers a user.
func (payload *User) Create() (map[string]interface{}, error) {

	// Check if there is any blank field
	var checkEmpty bool = (payload.Name == "" ||
		payload.Lastname == "" ||
		payload.Mail == "" ||
		payload.Password == "")
	if checkEmpty {
		return fiber.Map{
			"status":  false,
			"message": "Blank fields.",
		}, nil
	}

	// Validate first name
	if len(payload.Name) < 2 || len(payload.Name) > 32 {
		return fiber.Map{
			"status":  false,
			"message": "The name must be between 2 and 32 characters in length.",
		}, nil
	}

	// Validate last name
	if len(payload.Lastname) < 2 || len(payload.Lastname) > 32 {
		return fiber.Map{
			"status":  false,
			"message": "The lastname must be between 2 and 32 characters in length.",
		}, nil
	}

	// Validate password
	if isValid, reason := helpers.ValidatePassword(payload.Password); !isValid {
		return fiber.Map{
			"status":  false,
			"message": reason,
		}, nil
	}

	// Validate mail address
	_, err := mail.ParseAddress(payload.Mail)
	if err != nil {
		return fiber.Map{
			"status":  false,
			"message": "Invalid mail address.",
		}, nil
	}

	// Check if mail address exists in db
	exists, err := core.CheckExistsInDatabase(&User{Mail: payload.Mail})
	if err != nil {
		return nil, err
	}
	if exists {
		return fiber.Map{
			"status":  false,
			"message": "Unavailable mail address.",
		}, nil
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := User{
		Name:     payload.Name,
		Lastname: payload.Lastname,
		Mail:     payload.Mail,
		Password: hashedPassword,
	}

	err = core.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}

	return fiber.Map{
		"status":  true,
		"message": "OK!",
	}, nil

}

// Returns wallet if it exists, creates and returns it if doesnt exists.
func (user *User) GetOrCreateWallet(currency string) (*Wallet, error) {
	// Query payload
	wallet := Wallet{
		Currency: currency,
		UserID:   user.ID,
	}
	// Query execution
	result := core.DB.Model(&wallet).Where(&wallet).First(&wallet)
	if result.Error == nil {
		// If wallet found, return
		return &wallet, nil
	}

	// If wallet not found, create
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create wallet
		err := core.DB.Create(&wallet).Error
		if err != nil {
			// Failed to create wallet
			return nil, err
		}
		// Wallet successfully created, return it.
		return &wallet, nil
	}

	// If wallet search was failed, return error
	return nil, result.Error
}

// Returns wallet if it exists, returns nil if it doesnt.
func (user *User) GetWallet(currency string) (*Wallet, error) {
	// Query payload
	wallet := Wallet{
		Currency: currency,
		UserID:   user.ID,
	}
	// Query execution
	result := core.DB.Model(&wallet).Where(&wallet).First(&wallet)
	if result.Error == nil {
		// If wallet found, return
		return &wallet, nil
	}

	// If wallet search was failed, return error
	return nil, result.Error
}
