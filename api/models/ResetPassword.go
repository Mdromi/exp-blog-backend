package models

import (
	"gorm.io/gorm"
)

// ResetPassword represents a password reset token for a user
type ResetPassword struct {
	gorm.Model
	Email string `gorm:"size:100;not null;" json:"email"`
	Token string `gorm:"size:255;not null;" json:"token"`
}
