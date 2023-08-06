package models

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"gorm.io/gorm"
)

// Comment represents a comment on a post
type Comment struct {
	gorm.Model
	ProfileID uint32  `gorm:"not null" json:"profile_id"`
	PostID    uint64  `gorm:"not null" json:"post_id"`
	Body      string  `gorm:"type:text;not null" json:"body"`
	Profile   Profile `json:"profile"`
	// Replyes *Replyes `gorm:"foreignKey:CommentID" json:"replyes"`
}

func (c *Comment) Preapre() {
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.Profile = Profile{}
}

func (c *Comment) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if c.Body == "" {
			err = errors.New("required comment")
			errorMessages["Required_body"] = err.Error()
		}
	default:
		if c.Body == "" {
			err = errors.New("required comment")
			errorMessages["Required_body"] = err.Error()
		}
	}
	return errorMessages
}

func (c *Comment) SaveComment(db *gorm.DB) (*Comment, error) {
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&Profile{}).Where("id = ?", c.ProfileID).Take(&c.Profile).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) GetComments(db *gorm.DB, pid uint64) (*[]Comment, error) {
	comments := []Comment{}
	err := db.Debug().Model(&Comment{}).Where("post_id = ?", pid).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}

	// changed comments[i].UserID
	if len(comments) > 0 {
		for _, comment := range comments {
			err = db.Debug().Model(&Profile{}).Where("id = ?", comment.ProfileID).Take(&comment.Profile).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
		return &comments, err
	}
	return &comments, err
}

func (c *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {
	var err error

	err = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Updates(Comment{Body: c.Body}).Error

	if err != nil {
		return &Comment{}, err
	}

	fmt.Println("this is the comment body: ", c.Body)
	if c.ID != 0 {
		err = db.Debug().Model(&Profile{}).Where("id = ?", c.ProfileID).Take(&c.Profile).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

func (c *Comment) DeleteAComment(db *gorm.DB) (int64, error) {
	// Delete the comment replies first
	replyesModel := Replyes{}
	_, err := replyesModel.DeleteACommentReplyes(db, uint32(c.ID))
	if err != nil {
		return 0, err
	}

	// Now, delete the comment
	db = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// When a profile deleted, we also delete the comments that the profile had
func (c *Comment) DeleteUserComments(db *gorm.DB, profileID uint32) (int64, error) {
	commetns := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("profile_id", profileID).Find(&commetns).Delete(&commetns)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// When a post is deleted, we also delete the comments that the post had
func (c *Comment) DeletePostComments(db *gorm.DB, postID uint64) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("post_id = ?", postID).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
