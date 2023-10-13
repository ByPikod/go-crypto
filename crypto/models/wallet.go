package models

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	Currency    string        `json:"currency" gorm:"not null;index"`
	Balance     float64       `json:"balance" gorm:"default:0;not null"`
	UserID      uint          `json:"userID" gorm:"not null;index"`
	User        User          `gorm:"foreignKey:UserID"`
	Transaction []Transaction `gorm:"foreignKey:WalletID"`
}
