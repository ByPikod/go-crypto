package repositories

import (
	"errors"

	"github.com/ByPikod/go-crypto/tree/crypto/models"
	"gorm.io/gorm"
)

type (
	WalletRepository struct {
		db *gorm.DB
	}

	IWalletRepository interface {
		GetWallet(userID uint, currency string) (*models.Wallet, error)
		CreateWallet(wallet *models.Wallet) error
		SaveWallet(wallet *models.Wallet) error
		PreloadWallets(user *models.User) error
		CreateTransaction(transaction *models.Transaction) error
		RemoveTransactionByID(id uint) error
	}
)

var (
	ErrInvalidWallet          = errors.New("invalid wallet")
	ErrInvalidTransaction     = errors.New("invalid transaction")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

// Create new wallet repository
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// Returns wallet if exists, returns nil if not.
func (repository *WalletRepository) GetWallet(
	userID uint,
	currency string,
) (*models.Wallet, error) {

	// Query payload
	wallet := models.Wallet{
		Currency: currency,
		UserID:   userID,
	}

	// Query execution
	result := repository.db.Model(&wallet).Where(&wallet).First(&wallet)
	if result.Error == nil {
		// If wallet found, return it.
		return &wallet, nil
	}

	// Not found
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	// An unexpected error ocurred
	return nil, result.Error

}

// Creates a wallet.
func (repository *WalletRepository) CreateWallet(wallet *models.Wallet) error {
	return repository.db.Create(wallet).Error
}

// Save wallet data
func (repository *WalletRepository) SaveWallet(wallet *models.Wallet) error {
	if wallet.ID == 0 {
		return ErrInvalidWallet
	}
	return repository.db.Save(&wallet).Error

}

// Fetch all the wallets user has.
func (repository *WalletRepository) PreloadWallets(user *models.User) error {
	res := repository.db.Preload("Wallets").Find(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// Create a transaction
func (repository *WalletRepository) CreateTransaction(transaction *models.Transaction) error {

	// Create a new transaction
	result := repository.db.Create(&transaction)
	if result.Error != nil {
		return result.Error
	}

	// Successfully created
	return nil

}

// Remove a transaction by its ID
func (repository *WalletRepository) RemoveTransactionByID(id uint) error {

	payload := models.Transaction{}
	payload.ID = id

	// Remove the transaction
	result := repository.db.Delete(payload)
	if result.Error != nil {
		return result.Error
	}

	// Successfully removed
	return nil

}
