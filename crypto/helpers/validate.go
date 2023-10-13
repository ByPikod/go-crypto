package helpers

import (
	"net/mail"
	"strings"
)

func IsBlank(fields ...string) bool {
	for _, v := range fields {
		v = strings.TrimSpace(v)
		if v == "" {
			return true
		}
	}
	return false
}

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

// Validate user registration payload
func ValidateRegistration(
	name string,
	lastname string,
	mailAddress string,
	password string,
) (ok bool, reason string) {

	// Blank fiels
	var checkEmpty bool = IsBlank(name, lastname, mailAddress, password)
	if checkEmpty {
		return false, "Blank fields."
	}

	// First name
	if len(name) < 2 || len(name) > 32 {
		return false, "The name must be between 2 and 32 characters in length."
	}

	// Last name
	if len(lastname) < 2 || len(lastname) > 32 {
		return false, "The lastname must be between 2 and 32 characters in length."
	}

	// Password
	if isValid, reason := ValidatePassword(password); !isValid {
		return false, reason
	}

	// Mail address
	_, err := mail.ParseAddress(mailAddress)
	if err != nil {
		return false, "Invalid mail address."
	}

	return true, ""

}
