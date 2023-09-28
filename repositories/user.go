package repositories

import (
	"errors"

	"github.com/ByPikod/go-crypto/core"
	"github.com/ByPikod/go-crypto/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// Create user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Registers a user.
func (repo *UserRepository) Create(
	name string,
	lastname string,
	mailAddress string,
	password string,
) error {

	// Create user
	user := models.User{
		Name:     name,
		Lastname: lastname,
		Mail:     mailAddress,
		Password: password,
	}

	err := repo.db.Create(&user).Error
	return err

}

// Returns true if mail is available, false otherwise.
func (repo *UserRepository) CheckMailAvailable(mailAddress string) (bool, error) {
	exists, err := core.CheckExistsInDatabase(&models.User{Mail: mailAddress})
	if err != nil {
		return false, err
	}
	return !exists, nil
}

// Returns user by mail if found. Returns nil if could not found.
func (repo *UserRepository) GetUserByMail(mailAddress string) (*models.User, error) {
	row := models.User{Mail: mailAddress}
	res := core.DB.Model(&row).Where(&row).First(&row)
	if res.Error == nil {
		return &row, nil
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, res.Error
}
