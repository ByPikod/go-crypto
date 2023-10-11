package helpers

import (
	"errors"

	"gorm.io/gorm"
)

// Returns true if data exists in database.
func CheckExistsInDatabase(db *gorm.DB, dest interface{}, conds ...interface{}) (bool, error) {
	res := db.Model(dest).Where(dest, conds...).First(nil)
	if res.Error == nil {
		return true, nil // Query successfully executed and data found.
	}
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, nil // Error says data doesnt exists in db
	}
	return false, res.Error // We don't know if it exists in db or not but an error occured.
}