package repositories

import (
	"errors"

	"github.com/ByPikod/go-crypto/helpers"
	"github.com/ByPikod/go-crypto/models"
	"gorm.io/gorm"
)

type (
	UserRepository struct {
		db     *gorm.DB
		secret string
	}

	IUserRepository interface {
		Create(name string, lastname string, mailAddress string, password string) error
		IsMailAvailable(mailAddress string) (bool, error)
		GetUserByMail(mailAddress string) (*models.User, error)
		GetUserById(id uint) (*models.User, error)
		AuthSecret() string
	}
)

// Create user repository
func NewUserRepository(db *gorm.DB, secret string) *UserRepository {
	return &UserRepository{db: db}
}

// Returns the authentication secret
func (repo *UserRepository) AuthSecret() string {
	return repo.secret
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
func (repo *UserRepository) IsMailAvailable(mailAddress string) (bool, error) {
	exists, err := helpers.CheckExistsInDatabase(repo.db, &models.User{Mail: mailAddress})
	if err != nil {
		return false, err
	}
	return !exists, nil
}

// Returns user by mail if found. Returns nil if could not found.
func (repo *UserRepository) GetUserByMail(mailAddress string) (*models.User, error) {
	row := models.User{Mail: mailAddress}
	res := repo.db.Model(&row).Where(&row).First(&row)
	if res.Error == nil {
		return &row, nil
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, res.Error
}

// Returns user by its ID
func (repository *UserRepository) GetUserById(id uint) (*models.User, error) {

	// Paylaod
	ret := new(models.User)
	ret.ID = id

	// Query
	result := repository.db.Model(ret).Where(ret).First(ret)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return ret, nil

}
