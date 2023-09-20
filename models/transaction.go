package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Type     int8    `json:"type" gorm:"not null"`
	Change   float64 `json:"change" gorm:"not null"`
	Balance  float64 `json:"balance" gorm:"not null"`
	WalletID uint    `json:"walletID" gorm:"not null;index"`
}

// Why do we have a "change" field instead of directly the price?
// Simply, I don't want to make the change depend on the type of the transaction.
// If we have the price itself, then we need to make the calculations for each transaction type.
// And we will need to calculate the change for each transaction type.
// Thats why we have the "change" field instead of directly the price.
