package models

import (
	"errors"

	"github.com/ByPikod/go-crypto/core"
)

// Transaction types
const (
	TRANSACTION_TYPE_BUY      = -1
	TRANSACTION_TYPE_SELL     = 1
	TRANSACTION_TYPE_WITHDRAW = -2
	TRANSACTION_TYPE_DEPOSIT  = 2
)

// Error definitions
var (
	ErrInvalidTransaction     = errors.New("invalid transaction")
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

// Create a transaction
func (transaction *Transaction) Create() error {

	// Validate transaction type
	switch transaction.Type {
	case TRANSACTION_TYPE_BUY:
	case TRANSACTION_TYPE_SELL:
	case TRANSACTION_TYPE_WITHDRAW:
	case TRANSACTION_TYPE_DEPOSIT:
		break
	default:
		return ErrInvalidTransactionType
	}

	// Create a new transaction
	result := core.DB.Create(&transaction)
	if result.Error != nil {
		return result.Error
	}

	// Successfully created
	return nil

}

// Remove a transaction by its ID
func (transaction *Transaction) Remove() error {

	// Prepare payload
	if transaction.ID == 0 {
		// If ID is not set, occur an error
		return ErrInvalidTransaction
	}
	payload := Transaction{}
	payload.ID = transaction.ID

	// Remove the transaction
	result := core.DB.Delete(payload)
	if result.Error != nil {
		return result.Error
	}

	// Successfully removed
	return nil

}
