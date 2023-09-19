package models

import "gorm.io/gorm"

type TransactionHistory struct {
	gorm.Model
	Type         uint8   `json:"type" gorm:"not null"`
	Amount       float32 `json:"amount" gorm:"not null"`
	FinalBalance int     `json:"finalBalance" gorm:"not null"`
	Wallet       Wallet  `gorm:"not null;foreignKey:ID"`
}
