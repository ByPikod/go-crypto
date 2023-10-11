package core

import (
	"os"
	"strconv"

	"github.com/ByPikod/go-crypto/helpers"
	"github.com/joho/godotenv"
)

type (
	Configuration struct {
		Database               *DBInfo
		Loki                   *LokiInfo
		AuthSecret             string
		Host                   string
		Listen                 string
		ExchangesFetchInterval int
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
	// LokiInfo struct contains authentication information for the logging database.
	LokiInfo struct {
		Host string
	}
)

// It will be nil if config haven't been initialized.
var Config *Configuration

func init() {
	Config = InitializeConfig()
}

func or(x string, y string) string {
	if x == "" {
		return y
	}
	return x
}

// Calls the callback passed to it.
// If callback returns any error, function returns the default value.
// If callback is worked, returns the output of the callback.
//
// value:
//
//	Value to parse
//
// defaultValue:
//
//	Default value to return when parse failed.
//
// parse:
//
//	Callback function to parse.
func parse[T interface{}](
	value string,
	defaultValue T,
	parse func(value string) (T, error),
) T {
	res, err := parse(value)
	if err != nil {
		return defaultValue
	}
	return res
}

// Initializes config and makes Config variable above ready to use by loading environment variables.
// ".env" is supported.
func InitializeConfig() *Configuration {

	err := godotenv.Load()
	if err != nil {
		helpers.LogError(`File ".env" not found or cannot parsed: ` + err.Error())
	}

	db := &DBInfo{
		Port:     or(os.Getenv("DB_PORT"), "5432"),
		Host:     or(os.Getenv("DB_HOST"), "localhost"),
		User:     or(os.Getenv("DB_USER"), "postgres"),
		Password: or(os.Getenv("DB_PASS"), "root"),
		Database: or(os.Getenv("DB_NAME"), "gocrypto"),
	}

	loki := &LokiInfo{
		Host: or(os.Getenv("LOKI_HOST"), "http://loki:3100"),
	}

	config := Configuration{
		Database:   db,
		Loki:       loki,
		AuthSecret: or(os.Getenv("AUTH_SECRET"), "32f97916299787f211b5111e6da178b1"),
		Host:       or(os.Getenv("HOST"), ""),
		Listen:     or(os.Getenv("LISTEN"), "80"),
		ExchangesFetchInterval: parse[int](os.Getenv("ExchangesFetchInterval"), 30, func(value string) (int, error) {
			return strconv.Atoi(value)
		}),
	}

	return &config

}
