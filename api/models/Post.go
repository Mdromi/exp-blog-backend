package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title          string    `gorm:"size:255;not null" json:"title"`
	PostPermalinks string    `gorm:"size:255" json:"post_permalinks"`
	Body           string    `gorm:"type:text;not null" json:"body"`
	AuthorID       uint32    `gorm:"not null" json:"author_id"`
	Author         User      `gorm:"foreignKey:AuthorID" json:"author"`
	Tags           []string  `gorm:"type:text[]" json:"tags"`
	Thumbnails     string    `gorm:"size:255" json:"thumbnails"`
	ReadTime       string    `json:"read_time"`
	Likes          []Like    `gorm:"many2many:post_likes;" json:"likes"`
	DiLikes        []Disike  `gorm:"many2many:post_dislikes;" json:"dislikes"`
	Comment        []Comment `json:"comment"`
}
