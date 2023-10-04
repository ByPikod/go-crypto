package services

import (
	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/repositories"
	"github.com/gofiber/fiber/v2"
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

// Create new user service.
func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{repository: repository}
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
	token, err := helpers.GenerateUserToken(core.Config.AuthSecret, matching.ID)
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

func (service *UserService) GetUserById(id uint) (*models.User, error) {
	return service.repository.GetUserById(id)
}
