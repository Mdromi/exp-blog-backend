package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID     uint32  `gorm:"not null" json:"user_id"`
	Name       string  `gorm:"type:varchar(50);not null" json:"name" validate:"min=2,max=50"`
	Title      string  `gorm:"type:varchar(100);not null" json:"title" validate:"max=100"`
	Bio        string  `gorm:"type:text;not null" json:"bio" validate:"max=500"`
	ProfilePic string  `gorm:"type:varchar(255)" json:"profile_pic"`
	Links      Links   `json:"links"`
	Posts      []*Post `gorm:"many2many:profile_posts;" json:"posts"`
	Bookmarks  []*Post `gorm:"many2many:profile_bookmarks;" json:"bookmarks"`
	Flowing    []*User `gorm:"many2many:user_followers;association_foreignkey:FlowingID;" json:"flowing"`
}

type Links struct {
	gorm.Model
	Website  string `gorm:"type:varchar(255)" json:"website"`
	Facebook string `gorm:"type:varchar(255)" json:"facebook"`
	Twitter  string `gorm:"type:varchar(255)" json:"twitter"`
	Github   string `gorm:"type:varchar(255)" json:"github"`
}
