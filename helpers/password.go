package helpers

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Estimates the password strength and returns true if its enough. False otherwise.
func ValidatePassword(password string) (bool, string) {
	if !strings.ContainsAny(password, "1234567890") {
		return false, "Password must contain at least one number."
	}
	if !strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyzABCDEFGGHIJKLMNOPQRSTUVWXYZ") {
		return false, "Password must contain at least one letter."
	}
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long."
	}
	if len(password) > 255 {
		return false, "Password must be at most 255 characters long."
	}
	return true, ""
}

// Returns hashed version of the passed password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", nil
	}
	return string(bytes), err
}

// Compares two password hashes.
func ComparePasswords(hashedPassword string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}
