package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// LikeDislike model represents like/dislike actions
type LikeDislike struct {
	gorm.Model
	ProfileID uint   `gorm:"not null" json:"profile_id"`
	PostID    uint   `gorm:"not null" json:"post_id"`
	Action    string `gorm:"not null" json:"action"`
}

// Constants for different actions
const (
	ActionLike    = "like"
	ActionDislike = "dislike"
	ActionHard    = "hard"
	ActionSad     = "sad"
)

// SaveLikeDislike saves a like/dislike action to the database.
func (ld *LikeDislike) SaveLike(db *gorm.DB) (*LikeDislike, error) {
	// Check if the action is a valid one
	if !isValidAction(ld.Action) {
		return nil, errors.New("invalid action")
	}

	// Check if the user has previously disliked this post
	dislikedBefore, err := CheckIfDislikedBefore(db, ld.PostID, ld.ProfileID)
	if err != nil {
		return nil, err
	}

	if dislikedBefore != "" {
		// The user has previously disliked this post, so we need to remove the dislike
		err = RemoveLikeDislike(db, ld.PostID, ld.ProfileID)
		if err != nil {
			return nil, err
		}
		// The user has already performed this like/dislike action before, so return a custom error message
		if dislikedBefore == ld.Action {
			return nil, errors.New("you have already disliked this post")
		}
	}

	// The user has not performed this like/dislike action before, so let's save it
	newLikeDislike := &LikeDislike{
		ProfileID: ld.ProfileID,
		PostID:    ld.PostID,
		Action:    ld.Action,
	}
	err = db.Create(newLikeDislike).Error
	if err != nil {
		return nil, err
	}

	return newLikeDislike, nil
}

func (l *LikeDislike) DeleteLike(db *gorm.DB) (*LikeDislike, error) {
	var err error
	var deletedLike *LikeDislike

	err = db.Debug().Model(LikeDislike{}).Where("id = ?", l.ID).Take(&l).Error
	if err != nil {
		return &LikeDislike{}, err
	} else {
		// If the like exist, save it in deleted like and delete it
		deletedLike = l
		db = db.Debug().Model(&LikeDislike{}).Where("id = ?", l.ID).Take(&LikeDislike{}).Delete(&LikeDislike{})
		if db.Error != nil {
			fmt.Println("cant delete like", db.Error)
			return &LikeDislike{}, db.Error
		}
	}

	return deletedLike, nil
}

func (l *LikeDislike) GetLikesInfo(db *gorm.DB, pid uint) (*[]LikeDislike, error) {
	likeDislikes := []LikeDislike{}
	err := db.Debug().Model(&LikeDislike{}).Where("post_id = ?", pid).Find(&likeDislikes).Error
	if err != nil {
		return &[]LikeDislike{}, err
	}
	return &likeDislikes, err
}

// When a post is deleted, we also delete the likes that the post had
func (l *LikeDislike) DeleteUserLikes(db *gorm.DB, uid uint32) (int64, error) {
	likes := []LikeDislike{}
	db = db.Debug().Model(&LikeDislike{}).Where("profile_id = ?", uid).Find(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// When a post deleted, we also delete the likes that the post hat
func (l *LikeDislike) DeletePostLikes(db *gorm.DB, pid uint64) (int64, error) {
	likes := []LikeDislike{}
	db = db.Debug().Model(&LikeDislike{}).Where("post_id = ?", pid).Find(&likes).Delete(&likes)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// CheckIfDislikedBefore checks if the user has previously disliked a post.
func CheckIfDislikedBefore(db *gorm.DB, postID, profileID uint) (string, error) {
	likeDislike := LikeDislike{}

	err := db.Debug().Table("like_dislikes").Where("post_id = ? AND profile_id = ?", postID, profileID).Take(&likeDislike).Error
	if err == nil {
		return likeDislike.Action, nil // The user has previously disliked this post
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil // The user has not disliked this post before
	}
	return "", err
}

// RemoveDislike removes a dislike entry from the database for a post by a user.
func RemoveLikeDislike(db *gorm.DB, postID, profileID uint) error {
	likeDislike := LikeDislike{}
	return db.Debug().Delete(&likeDislike, "post_id = ? AND profile_id = ?", postID, profileID).Error
}

func isValidAction(action string) bool {
	switch action {
	case ActionLike, ActionDislike, ActionHard, ActionSad:
		return true
	default:
		return false
	}
}
