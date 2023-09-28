package services

import "github.com/ByPikod/go-crypto/repositories"

type WalletService struct {
	repository *repositories.WalletRepository
}

// Create new wallet service
func NewWalletService(repository *repositories.WalletRepository) *WalletService {
	return &WalletService{repository: repository}
}
