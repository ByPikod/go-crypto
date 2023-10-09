package core

import (
	"os"

	"github.com/ByPikod/go-crypto/helpers"
	"github.com/joho/godotenv"
)

type (
	Configuration struct {
		AuthSecret string
		Database   *DBInfo
		Host       string
		Listen     string
	}
	// DBInfo struct contains authentication information for the database.
	// Used to configure the database. See core/config
	DBInfo struct {
		Host     string
		Port     string
		User     string
		Password string
		Database string
	}
)

// It will be nil if config haven't been initialized.
var Config *Configuration

func or(x string, y string) string {
	if x == "" {
		return y
	}
	return x
}

// Initializes config and makes Config variable above ready to use by loading environment variables.
// ".env" is supported.
func InitializeConfig() *Configuration {

	err := godotenv.Load()
	if err != nil {
		helpers.LogError(`File ".env" not found or cannot parsed.`)
	}

	dbInfo := DBInfo{
		Port:     or(os.Getenv("DB_PORT"), "5432"),
		Host:     or(os.Getenv("DB_HOST"), "localhost"),
		User:     or(os.Getenv("DB_USER"), "postgres"),
		Password: or(os.Getenv("DB_PASS"), "root"),
		Database: or(os.Getenv("DB_NAME"), "gocrypto"),
	}

	config := Configuration{
		AuthSecret: or(os.Getenv("AUTH_SECRET"), "32f97916299787f211b5111e6da178b1"),
		Database:   &dbInfo,
		Host:       or(os.Getenv("HOST"), ""),
		Listen:     or(os.Getenv("LISTEN"), "8080"),
	}

	Config = &config
	return &config

}
