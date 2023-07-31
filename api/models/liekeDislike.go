package models

import (
	"gorm.io/gorm"
)

type LikeDislike struct {
	gorm.Model
	ProfileID uint   `gorm:"not null" json:"profile_id"`
	PostID    uint   `gorm:"not null" json:"post_id"`
	Action    string `gorm:"not null" json:"action"`
}

const (
	ActionLike    = "like"
	ActionDislike = "dislike"
	ActionHard    = "hard"
	ActionSad     = "sad"
)

// CheckIfDislikedBefore checks if the user has previously disliked a post.
func CheckIfDislikedBefore(db *gorm.DB, postID, profileID uint) (bool, error) {
	dislike := Dislike{}

	err := db.Debug().Model(&Dislike{}).Where("post_id = ? AND profile_id = ?", postID, profileID).Take(&dislike).Error
	if err == nil {
		return true, nil // The user has previously disliked this post
	} else if err == gorm.ErrRecordNotFound {
		return false, nil // The user has not disliked this post before
	}
	return false, err
}

// RemoveDislike removes a dislike entry from the database for a post by a user.
func RemoveDislike(db *gorm.DB, postID, profileID uint) error {
	dislike := Dislike{}
	return db.Debug().Delete(&dislike, "post_id = ? AND profile_id = ?", postID, profileID).Error
}

// CheckIfLikedBefore checks if the user has previously liked a post.
func CheckIfLikedBefore(db *gorm.DB, postID, profileID uint) (bool, error) {
	like := Like{}
	err := db.Debug().Model(&Like{}).Where("post_id = ? AND profile_id = ?", postID, profileID).Take(&like).Error
	if err == nil {
		return true, nil // The user has liked this post before
	} else if err == gorm.ErrRecordNotFound {
		return false, nil // The user has not liked this post before
	}
	return false, err
}

// SaveNewLike saves a new like entry to the database for a post by a user.
func SaveNewLike(db *gorm.DB, postID, profileID uint) error {
	like := &Like{
		ProfileID: profileID,
		PostID:    postID,
	}
	return db.Debug().Create(like).Error
}

// RemoveLike removes a like entry from the database for a post by a user.
func RemoveLike(db *gorm.DB, postID, profileID uint) error {
	like := Like{}
	return db.Debug().Delete(&like, "post_id = ? AND profile_id = ?", postID, profileID).Error
}

// SaveNewDislike saves a new dislike entry to the database for a post by a user.
func SaveNewDislike(db *gorm.DB, postID, profileID uint) error {
	dislike := &Dislike{
		ProfileID: profileID,
		PostID:    postID,
	}
	return db.Debug().Create(dislike).Error
}
