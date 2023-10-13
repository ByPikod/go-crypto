package repositories

import (
	"errors"

	"github.com/ByPikod/go-crypto/tree/crypto/core"
	"github.com/ByPikod/go-crypto/tree/crypto/helpers"
	"github.com/ByPikod/go-crypto/tree/crypto/models"
	"gorm.io/gorm"
)

type (
	UserRepository struct {
		db *gorm.DB
	}

	IUserRepository interface {
		AuthSecret() string
		CreateUser(*models.User) error
		SaveUser(user *models.User) error
		IsMailAvailable(mailAddress string) (bool, error)
		GetUserByMail(mailAddress string) (*models.User, error)
		GetUserById(id uint) (*models.User, error)
		CreateVerification(verification *models.Verification) error
		UpdateVerification(verification *models.Verification) error
		GetVerification(verification models.Verification) (*models.Verification, error)
		RemoveVerificationById(ID uint) error
	}
)

// Create user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Returns the authentication secret
func (repo *UserRepository) AuthSecret() string {
	return core.Config.AuthSecret
}

// Registers a user.
func (repo *UserRepository) CreateUser(user *models.User) error {
	return repo.db.Create(&user).Error
}

func (repo *UserRepository) SaveUser(user *models.User) error {
	return repo.db.Save(&user).Error
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

// Returns user by its ID, returns nil if user not found.
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

// Create verification
func (repository *UserRepository) CreateVerification(verification *models.Verification) error {
	return repository.db.Create(verification).Error
}

// Update verification
func (repository *UserRepository) UpdateVerification(verification *models.Verification) error {
	return repository.db.Save(verification).Error
}

// Returns verification record by its ID, returns nil if not found.
func (repository *UserRepository) GetVerification(verification models.Verification) (*models.Verification, error) {

	// Query
	result := repository.db.Model(&verification).Where(&verification).First(&verification)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &verification, nil

}

func (repository *UserRepository) RemoveVerificationById(ID uint) error {
	verification := new(models.Verification)
	verification.ID = ID
	return repository.db.Delete(verification).Error
}
