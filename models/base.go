package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"createdAt,omitempty" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updatedAt,omitempty" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}
