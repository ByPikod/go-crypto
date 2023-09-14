package core

import (
	"os"

	"github.com/ByPikod/go-crypto/helpers"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Database *DBInfo
}

// DBInfo struct contains authentication information for the database.
// Used to configure the database. See core/config
type DBInfo struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// It will be nil if config haven't been initialized.
var Config *Configuration

// Initializes config and makes Config variable above ready to use by loading environment variables.
// ".env" is supported.
func InitializeConfig() {

	err := godotenv.Load()
	if err != nil {
		helpers.LogError(`An error occurred while loading environment variables.`)
		panic(err)
	}

	dbInfo := DBInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Database: os.Getenv("DB_NAME"),
	}

	config := Configuration{
		Database: &dbInfo,
	}

	Config = &config

}
