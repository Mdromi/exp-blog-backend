package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserID  uint32     `gorm:"not null" json:"user_id"`
	PostID  uint64     `gorm:"not null" json:"post_id"`
	Body    string     `gorm:"type:text;not null" json:"body"`
	User    User       `gorm:"foreignKey:UserID" json:"user"`
	Replyes []*Replyes `gorm:"foreignKey:CommentID" json:"replyes"`
}
