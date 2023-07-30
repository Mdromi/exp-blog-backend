package models

import (
	"gorm.io/gorm"
)

// Replyes represents a reply to a comment
type Replyes struct {
	gorm.Model
	CommentID uint64 `gorm:"not null" json:"comment_id"`
	UserID    uint32 `gorm:"not null" json:"user_id"`
	PostID    uint64 `gorm:"not null" json:"post_id"`
	Body      string `gorm:"type:text;not null" json:"body"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}
