package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Currency string `json:"currency" gorm:"not null;index"`
	Balance  int    `json:"balance" gorm:"default:0;not null"`
	User     User   `gorm:"not null;foreignKey:ID"`
}
