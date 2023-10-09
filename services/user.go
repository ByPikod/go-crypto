package services

import (
	"errors"
	"fmt"

	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type (
	UserService struct {
		repository repositories.IUserRepository
	}
	RegisterPaylaod struct {
		Name     string `json:"name"`
		Lastname string `json:"lastName"`
		Mail     string `json:"mail"`
		Password string `json:"password"`
	}
)

var (
	ErrTokenMalformed     = errors.New("token malformed")
	ErrTokenCouldntParsed = errors.New("token signed but id is not uint")
	ErrTokenUnauthorized  = errors.New("unauthorized")
)

// Create new user service.
func NewUserService(repository repositories.IUserRepository) *UserService {
	return &UserService{repository: repository}
}

// Authenticates a user from its token.
//
// It will return user id, if successfully authenticated.
// Will occur an error if not.
func (service *UserService) Authenticate(token string) (uint, error) {

	// Token parse
	tokenObj, err := jwt.Parse(token, func(tokenObj *jwt.Token) (interface{}, error) {
		return []byte(service.repository.AuthSecret()), nil
	})

	if err != nil {
		// Failed to parse token
		fmt.Println("failed to parse token.")
		return 0, ErrTokenMalformed
	}

	// Get claims by decoding the token
	claims, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("failed to parse claims.")
		return 0, ErrTokenUnauthorized
	}

	userID, ok := claims["UserID"].(float64)
	if !ok {
		return 0, ErrTokenCouldntParsed
	}

	return uint(userID), nil

}

// Generates a token for a user by its ID.
func (service *UserService) GenerateUserToken(id uint) (string, error) {

	claims := struct {
		jwt.StandardClaims
		UserID uint
	}{
		UserID: id,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &claims)
	tokenString, err := token.SignedString([]byte(service.repository.AuthSecret()))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

// Registers a user.
func (service *UserService) Create(
	name string,
	lastname string,
	mailAddress string,
	password string,
) (map[string]interface{}, error) {

	// Validate data
	ok, reason := helpers.ValidateRegistration(
		name,
		lastname,
		mailAddress,
		password,
	)

	if !ok {
		return fiber.Map{
			"status":  ok,
			"message": reason,
		}, nil
	}

	// Check if mail address available in db
	available, err := service.repository.IsMailAvailable(mailAddress)
	if err != nil {
		return nil, err
	}

	if !available {
		return fiber.Map{
			"status":  false,
			"message": "Unavailable mail address.",
		}, nil
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	service.repository.Create(
		name,
		lastname,
		mailAddress,
		hashedPassword,
	)

	return fiber.Map{
		"status":  true,
		"message": "OK!",
	}, nil

}

func (service *UserService) Login(mailAddress string, password string) (map[string]interface{}, error) {

	const incorrectCredentials = "Incorrect credentials!"
	// Check blank fields
	if helpers.IsBlank(mailAddress, password) {
		return fiber.Map{
			"status":  false,
			"message": "Blank fields.",
		}, nil
	}

	// Find the matching mail address from the database.
	matching, err := service.repository.GetUserByMail(mailAddress)

	if err != nil {
		return nil, err
	}

	if matching == nil {
		return fiber.Map{
			"status":  false,
			"message": incorrectCredentials,
		}, nil
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

	// Create token
	token, err := service.GenerateUserToken(matching.ID)
	if err != nil {
		return nil, err
	}

	// Return token
	return fiber.Map{
		"status":  true,
		"message": "OK!",
		"token":   token,
	}, nil

}

// Get user by its ID
func (service *UserService) GetUserById(id uint) (*models.User, error) {
	return service.repository.GetUserById(id)
}
