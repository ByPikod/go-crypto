package services

import (
	"github.com/ByPikod/go-crypto/models"
	"github.com/ByPikod/go-crypto/repositories"
)

type WalletService struct {
	repository repositories.IWalletRepository
}

// Create new wallet service
func NewWalletService(repository repositories.IWalletRepository) *WalletService {
	return &WalletService{repository: repository}
}

// Returns wallet if exists. Creates it and returns it if not exists.
func (service *WalletService) GetOrCreateWallet(userID uint, currency string) (*models.Wallet, error) {

	// Try to get wallet
	result, err := service.repository.GetWallet(userID, currency)
	if err != nil {
		return nil, err
	}

	if result != nil {
		return result, nil
	}

	// Payload
	wallet := &models.Wallet{
		UserID:   userID,
		Currency: currency,
	}

	// Create wallet
	err = service.repository.CreateWallet(wallet)
	if err != nil {
		return nil, err
	}

	return wallet, err

}

// Returns wallet if it exists, returns nil if it doesnt.
func (service *WalletService) GetWallet(userID uint, currency string) (*models.Wallet, error) {
	return service.repository.GetWallet(userID, currency)
}

// Loads wallets into user.
func (service *WalletService) LoadWallets(user *models.User) error {
	return service.repository.PreloadWallets(user)
}

// Save wallet data
func (service *WalletService) SaveWallet(wallet *models.Wallet) error {
	return service.repository.SaveWallet(wallet)
}

// Add transaction to wallet
func (service *WalletService) AddTransaction(wallet *models.Wallet, transactionType int8, change float64) (*models.Transaction, error) {

	// Wallet validation
	if wallet.ID == 0 {
		return nil, repositories.ErrInvalidWallet
	}

	// Transaction type validation
	switch transactionType {
	case models.TRANSACTION_TYPE_BUY:
	case models.TRANSACTION_TYPE_SELL:
	case models.TRANSACTION_TYPE_WITHDRAW:
	case models.TRANSACTION_TYPE_DEPOSIT:
		break
	default:
		return nil, repositories.ErrInvalidTransactionType
	}
	wallet.Balance += change

	// Transaction payload
	transaction := &models.Transaction{
		Type:     transactionType,
		Change:   change,
		Balance:  wallet.Balance,
		WalletID: wallet.ID,
	}

	// Create transaction
	err := service.repository.CreateTransaction(transaction)
	if err != nil {
		// Failed to create transaction
		return nil, err
	}

	// Save wallet
	err = service.repository.SaveWallet(wallet)
	if err != nil {
		// Failed to save wallet
		removeErr := service.repository.RemoveTransactionByID(transaction.ID)
		if removeErr != nil {
			panic(removeErr)
		}
		return nil, err
	}

	return transaction, nil
}
