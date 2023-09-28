package models

import (
	"errors"

	"github.com/ByPikod/go-crypto/core"
	"gorm.io/gorm"
)

/*
 * Functions
 */

// Retrieves user by its ID, returns nil if user not found.
func GetUserById(id uint) (*User, error) {
	ret := User{}
	ret.ID = id
	result := core.DB.Model(&ret).Where(&ret).First(&ret)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &ret, nil
}

// Returns wallet if it exists, creates and returns it if doesnt exists.
func (user User) GetOrCreateWallet(currency string) (*Wallet, error) {

	// Query payload
	wallet := Wallet{
		Currency: currency,
		UserID:   user.ID,
	}
	// Query execution
	result := core.DB.Model(&wallet).Where(&wallet).First(&wallet)
	if result.Error == nil {
		// If wallet found, return
		return &wallet, nil
	}

	// If wallet not found, create
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Create wallet
		err := core.DB.Create(&wallet).Error
		if err != nil {
			// Failed to create wallet
			return nil, err
		}
		// Wallet successfully created, return it.
		return &wallet, nil
	}

	// If wallet search was failed, return error
	return nil, result.Error

}

// Returns wallet if it exists, returns nil if it doesnt.
func (user User) GetWallet(currency string) (*Wallet, error) {

	// Query payload
	wallet := Wallet{
		Currency: currency,
		UserID:   user.ID,
	}
	// Query execution
	result := core.DB.Model(&wallet).Where(&wallet).First(&wallet)
	if result.Error == nil {
		// If wallet found, return
		return &wallet, nil
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	// If wallet search was failed, return error
	return nil, result.Error

}

// Fetch all the wallets user has.
func (user *User) PreloadWallets() error {
	res := core.DB.Preload("Wallets").Find(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
