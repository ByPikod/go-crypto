package helpers

import (
	"math/rand"
	"time"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand

// Generate seed
func init() {
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Algorithm for randomly generating verification code.
func GenerateVerificationCode(length int) string {
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(code)
}
