package models

import (
	"gorm.io/gorm"
)

type Replyes struct {
	gorm.Model
	CommentID uint64 `gorm:"not null" json:"comment_id"`
	UserID    uint32 `gorm:"not null" json:"user_id"`
	PostID    uint64 `gorm:"not null" json:"post_id"`
	Body      string `gorm:"type:text;not null" json:"body"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
	// Comments  []Comment `gorm:"many2many:reply_comments;" json:"comments"`
}
