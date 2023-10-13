package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ByPikod/go-crypto/tree/crypto/core"
	"github.com/ByPikod/go-crypto/tree/crypto/helpers"
	"github.com/ByPikod/go-crypto/tree/crypto/models"
	"github.com/ByPikod/go-crypto/tree/crypto/repositories"
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
	ErrTokenMalformed         = errors.New("token malformed")
	ErrTokenCouldntParsed     = errors.New("token signed but id is not uint")
	ErrTokenUnauthorized      = errors.New("unauthorized")
	ErrVerificationCooldown   = errors.New("wait to generate a new code")
	ErrVerificationNeeded     = errors.New("verification code needed")
	ErrInvalidVerification    = errors.New("invalid verification code")
	ErrUnavailableMailAddress = errors.New("mail address unavailable")
	ErrVerificationBlocked    = errors.New("email verification blocked due to too many failed attempts")
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
	verification string,
) error {

	// Check if mail address available in db
	available, err := service.repository.IsMailAvailable(mailAddress)
	if err != nil {
		// Internal server error
		return err
	}

	if !available {
		// Mail is unavailable
		return ErrUnavailableMailAddress
	}

	// Check mail verification
	verification = strings.TrimSpace(verification)
	if len(verification) == 0 {
		// Verification is not passed.
		return ErrVerificationNeeded
	}

	// Check if verification code is valid
	verificationData, err := service.repository.GetVerification(models.Verification{
		Mail: mailAddress,
	})

	if err != nil {
		// An error ocurred while requesting for verification data.
		return err
	}

	if verificationData == nil {
		// There is no code created for this mail address.
		return ErrVerificationNeeded
	}

	if verificationData.Fails > 10 {
		// Verification blocked due to failed attempts.
		return ErrVerificationBlocked
	}

	if verificationData.Code != strings.ToUpper(verification) {
		// Invalid verification code
		verificationData.Fails += 1
		service.repository.UpdateVerification(verificationData)
		return ErrInvalidVerification
	}

	// Hash password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}

	// Create user
	service.repository.CreateUser(&models.User{
		Name:     name,
		Lastname: lastname,
		Mail:     mailAddress,
		Password: hashedPassword,
	})

	// Successfully registered
	return nil

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

// Creates a new verification or updates the existing one.
//
// Will occur an error (ErrVerificationCooldown) if last verification mail is just created.
// Will return the code if successfully created.
func (service *UserService) CreateNewVerification(mail string) (string, error) {

	// Query
	verification, err := service.repository.GetVerification(models.Verification{
		Mail: mail,
	})

	if err != nil {
		return "", err
	}

	// Generate a random code for putting it to database.
	verificationCode := helpers.GenerateVerificationCode(6)

	// Create a new verification code
	if verification == nil {
		err = service.repository.CreateVerification(&models.Verification{
			Mail:  mail,
			Code:  verificationCode,
			Fails: 0,
		})

		if err != nil {
			// An error ocurred while inserting code to the database
			return "", err
		}

		// Return the code
		return verificationCode, nil
	}

	// Update the existing verification code
	now := time.Now()
	cooldown := (time.Second * core.Config.VerificationCooldown)
	if now.Sub(verification.UpdatedAt) < cooldown {
		// Feature is in cooldown.
		return "", ErrVerificationCooldown
	}

	verification.Code = verificationCode
	verification.Fails = 0
	err = service.repository.UpdateVerification(verification)
	if err != nil {
		return "", err
	}

	return verificationCode, nil

}

// Get user by its ID
func (service *UserService) GetUserById(id uint) (*models.User, error) {
	return service.repository.GetUserById(id)
}
