package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string  `gorm:"size:255;not null;unique" json:"username" validate:"min=2"`
	Email         string  `gorm:"size:100;not null;unique" json:"email"`
	Password      string  `gorm:"size:100;not null;" json:"password"`
	ProfileID     uint    `json:"profile_id"`
	Profile       Profile `json:"profile"`
	AvatarPath    string  `gorm:"size:255;" json:"avatar_path"`
	LikedPosts    []*Post `gorm:"many2many:user_liked_posts;" json:"liked_posts"`
	DislikedPosts []*Post `gorm:"many2many:user_disliked_posts;" json:"disliked_posts"`
}
