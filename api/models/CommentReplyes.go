package models

import (
	"errors"
	"fmt"
	"html"
	"strings"

	"gorm.io/gorm"
)

// Replyes represents a reply to a comment
type Replyes struct {
	gorm.Model
	CommentID uint64  `gorm:"not null" json:"comment_id"`
	PostID    uint32  `gorm:"not null" json:"post_id"`
	ProfileID uint64  `gorm:"not null" json:"profile_id"`
	Body      string  `gorm:"type:text;not null" json:"body"`
	Profile   Profile `json:"profile"`
}

func (rc *Replyes) Preapre() {
	rc.Body = html.EscapeString(strings.TrimSpace(rc.Body))
	rc.Profile = Profile{}
}

func (rc *Replyes) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if rc.Body == "" {
			err = errors.New("required comment")
			errorMessages["Required_body"] = err.Error()
		}
	default:
		if rc.Body == "" {
			err = errors.New("required comment")
			errorMessages["Required_body"] = err.Error()
		}
	}
	return errorMessages
}

func (rc *Replyes) SaveCommentReplyes(db *gorm.DB) (*Replyes, error) {
	// Check if the comment exists with the given CommentID
	var comment Comment
	if err := db.Debug().First(&comment, rc.CommentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("comment with ID %d not found", rc.CommentID)
		}
		return nil, err
	}

	// Save the comment reply
	if err := db.Debug().Create(&rc).Error; err != nil {
		return nil, err
	}

	return rc, nil
}

func (rc *Replyes) GetCommentReplyes(db *gorm.DB, cid uint64) (*[]Replyes, error) {
	replyes := []Replyes{}
	err := db.Debug().Model(&Replyes{}).Where("comment_id = ?", cid).Order("created_at desc").Find(&replyes).Error

	if err != nil {
		return &[]Replyes{}, err
	}

	fmt.Println("this is the comment body: ", rc.Body)
	if rc.ID != 0 {
		err = db.Debug().Model(&Profile{}).Where("id =?", rc.ProfileID).Take(&rc.Profile).Error
		if err != nil {
			return &[]Replyes{}, err
		}
	}
	return &replyes, nil
}

func (rc *Replyes) UpdateACommentReplyes(db *gorm.DB) (*Replyes, error) {
	var err error

	// Check if the comment exists with the given CommentID
	// TASK: This part we can deleted, cz we same validation controller function
	var comment Comment
	if err := db.Debug().First(&comment, rc.CommentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("comment with ID %d not found", rc.CommentID)
		}
		return nil, err
	}

	err = db.Debug().Model(&Replyes{}).Where("id = ?", rc.CommentID).Updates(Comment{Body: rc.Body}).Error
	if err != nil {
		return &Replyes{}, err
	}

	fmt.Println("this is the comment body: ", rc.Body)
	if rc.ID != 0 {
		err = db.Debug().Model(&Profile{}).Where("id =?", rc.ProfileID).Take(&rc.Profile).Error
		if err != nil {
			return &Replyes{}, err
		}
	}
	return rc, nil
}

func (rc *Replyes) DeleteAReplyes(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Replyes{}).Where("id = ?", rc.ID).Take(&Replyes{}).Delete(&Replyes{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// When a comment deleted, we also delete the comments replyes that the comments had
func (rc *Replyes) DeleteACommentReplyes(db *gorm.DB, commentID uint32) (int64, error) {
	// Check if there are any replies associated with the given commentID
	var count int64
	if err := db.Debug().Model(&Replyes{}).Where("comment_id = ?", commentID).Count(&count).Error; err != nil {
		return 0, err
	}

	if count == 0 {
		// There are no replies to delete, return 0 rows affected.
		return 0, nil
	}

	// Delete the replies
	db = db.Debug().Model(&Replyes{}).Where("comment_id = ?", commentID).Delete(&Replyes{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// When a profile deleted, we also delete the comment replyes that the profile had
func (rc *Replyes) DeleteUserProfileCommentReplyes(db *gorm.DB, profileID uint32) (int64, error) {
	replyes := []Replyes{}
	db = db.Debug().Model(&Replyes{}).Where("profile_id", profileID).Find(&replyes).Delete(&replyes)

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// When a post is deleted, we also delete the comment replyes that the post had
func (c *Replyes) DeletePostCommentReplyes(db *gorm.DB, postID uint64) (int64, error) {
	replyes := []Replyes{}
	db = db.Debug().Model(&Replyes{}).Where("post_id = ?", postID).Find(&replyes).Delete(&replyes)

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
