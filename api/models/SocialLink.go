package models

import (
	"gorm.io/gorm"
)

// SocialLink represents social links for a user's profile
type SocialLink struct {
	gorm.Model
	ProfileID uint32 `gorm:"not null" json:"profile_id"`
	Website   string `gorm:"type:varchar(255)" json:"website"`
	Facebook  string `gorm:"type:varchar(255)" json:"facebook"`
	Twitter   string `gorm:"type:varchar(255)" json:"twitter"`
	Github    string `gorm:"type:varchar(255)" json:"github"`
}
