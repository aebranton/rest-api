package database

import (
	"github.com/aebranton/rest-api/internal/user"
	"github.com/jinzhu/gorm"
)

// MigrateDB - migrates our database, creating the user table and cols
func MigrateDB(db *gorm.DB) error {
	if result := db.AutoMigrate(&user.User{}); result.Error != nil {
		return result.Error
	}
	return nil
}
