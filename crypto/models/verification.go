package models

import "gorm.io/gorm"

type (
	Verification struct {
		gorm.Model
		Mail  string `json:"mail" gorm:"index;not null;unique"`
		Code  string `json:"code" gorm:"not null"`
		Fails uint8  `json:"fails" gorm:"default:0"`
	}
)
