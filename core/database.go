package core

import (
	"fmt"

	"github.com/ByPikod/go-crypto/tree/crypto/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initializes the database
// Pointer receiver is the DBInfo object and can be accessed via core.Config.Database
func InitializeDatabase(dbInfo *DBInfo) *gorm.DB {
	// Database configuration
	db, err := dbInfo.connect()
	if err != nil {
		panic(err)
	}

	return db
}

// Creates a connection to the database and returns the Database object.
// If the database is not existing then creates a new one and connects to it.
// This is a private method and called from IntializeDatabase()
func (dbInfo *DBInfo) connect() (*gorm.DB, error) {

	// Connection without database name
	connString := fmt.Sprintf(
		"host=%s user=%s password=%s port=%s sslmode=disable",
		dbInfo.Host,
		dbInfo.User,
		dbInfo.Password,
		dbInfo.Port,
	)
	// Connection with database name
	connStringWithDatabase := fmt.Sprintf("%s dbname=%s", connString, dbInfo.Database)

	// Connect to Postgres without db name and create database if not exists. And connect to it.
	helpers.LogInfo("Connecting to the postgres server.")
	dbConn, err := gorm.Open(postgres.Open(connString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	helpers.LogInfo("Successfully connected to the postgres server.")

	dbFound := 0
	dbConn.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", dbInfo.Database).Scan(&dbFound)

	if dbFound == 0 {
		helpers.LogInfo("Database not found, creating new database...")
		result := dbConn.Exec("CREATE DATABASE " + dbInfo.Database)

		if result.Error != nil {
			return nil, result.Error
		}

		helpers.LogInfo("Database successfully created.")
	}

	// Close temporary database connection
	sqlDB, err := dbConn.DB()
	if err != nil {
		helpers.LogError("Failed to close temporary database connection: " + err.Error())
		return nil, err
	}
	sqlDB.Close()

	// Create actual connection
	db, err := gorm.Open(postgres.Open(connStringWithDatabase), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	helpers.LogError("Established database connection.")
	return db, nil
}
