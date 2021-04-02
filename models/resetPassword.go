package models

import(
	_"github.com/jinzhu/gorm"
	"time"
)
type ResetPassword struct {
	Email      string `gorm:"not null;primary_key"`
	Code       string `gorm:"not null"`
	CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"update_at"`
}

