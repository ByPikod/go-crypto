package models

import (
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Lastname string `json:"lastName" gorm:"not null"`
	Mail     string `json:"mail" gorm:"index;not null"`
	Password string `json:"password" gorm:"not null"`
}
