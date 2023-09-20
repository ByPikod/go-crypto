package models

import (
	"errors"

	"github.com/ByPikod/go-crypto/core"
)

var (
	ErrInvalidWallet = errors.New("invalid wallet")
)

// Save wallet data
func (wallet *Wallet) Save() error {
	if wallet.ID == 0 {
		return ErrInvalidWallet
	}
	return core.DB.Save(&wallet).Error
}

// Add transaction to wallet
func (wallet *Wallet) AddTransaction(transactionType int8, change float64) (*Transaction, error) {
	// Wallet validation
	if wallet.ID == 0 {
		return nil, ErrInvalidWallet
	}

	// Transaction type validation
	switch transactionType {
	case TRANSACTION_TYPE_BUY:
	case TRANSACTION_TYPE_SELL:
	case TRANSACTION_TYPE_WITHDRAW:
	case TRANSACTION_TYPE_DEPOSIT:
		break
	default:
		return nil, ErrInvalidTransactionType
	}
	wallet.Balance += change

	// Transaction payload
	transaction := &Transaction{
		Type:     transactionType,
		Change:   change,
		Balance:  wallet.Balance,
		WalletID: wallet.ID,
	}

	// Create transaction
	err := transaction.Create()
	if err != nil {
		// Failed to create transaction
		return nil, err
	}

	// Save wallet
	err = wallet.Save()
	if err != nil {
		// Failed to save wallet
		removeErr := transaction.Remove()
		if removeErr != nil {
			panic(removeErr)
		}
		return nil, err
	}

	return transaction, nil
}
