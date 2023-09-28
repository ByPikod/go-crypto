package repositories

import "gorm.io/gorm"

type WalletRepository struct {
	db *gorm.DB
}

// Create new wallet repository
func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}
