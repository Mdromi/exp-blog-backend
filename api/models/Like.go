package models

import (
	"errors"

	"gorm.io/gorm"
)

// Like model represents a like by a user on a post
type Like struct {
	gorm.Model
	ProfileID uint `gorm:"not null" json:"profile_id"`
	PostID    uint `gorm:"not null" json:"post_id"`
}

func (l *Like) SaveLike(db *gorm.DB) (*Like, error) {
	// Check if the user has previously disliked this post
	dislikedBefore, err := CheckIfDislikedBefore(db, l.PostID, l.ProfileID)
	if err != nil {
		return &Like{}, err
	}

	if dislikedBefore {
		// The user has previously disliked this post, so we need to remove the dislike
		err = RemoveDislike(db, l.PostID, l.ProfileID)
		if err != nil {
			return &Like{}, err
		}

		// Update the profile's DislikedPosts field to remove the post
		// err = RemoveFromDislikedPosts(db, l.PostID, l.ProfileID)
		// if err != nil {
		// 	return &Like{}, err
		// }
	}

	// Check if the user has liked this post before
	likedBefore, err := CheckIfLikedBefore(db, l.PostID, l.ProfileID)
	if err != nil {
		return &Like{}, err
	}

	if !likedBefore {
		// The user has not liked this post before, so let's save the incoming like
		err = SaveNewLike(db, l.PostID, l.ProfileID)
		if err != nil {
			return &Like{}, err
		}

	} else {
		// The user has liked it before, so return a custom error message
		return &Like{}, errors.New("double like")
	}

	return l, nil
}

func (l *Like) DeleteLike(db *gorm.DB) (*Like, error) {
	// Check if the user has liked this post before
	likedBefore, err := CheckIfLikedBefore(db, l.PostID, l.ProfileID)
	if err != nil {
		return &Like{}, err
	}

	if likedBefore {
		// The user has liked this post, so we need to remove the like
		err = RemoveLike(db, l.PostID, l.ProfileID)
		if err != nil {
			return &Like{}, err
		}

		// Update the profile's LikedPosts field to remove the post
		// err = RemoveFromLikedPosts(db, l.PostID, l.ProfileID)
		// if err != nil {
		// 	return &Like{}, err
		// }
	} else {
		// The user has not liked this post before, so return a custom error message
		return &Like{}, errors.New("like not found")
	}

	return l, nil
}

func (l *Like) GetLikesInfo(db *gorm.DB, pid uint64) (*[]Like, error) {
	likes := []Like{}
	err := db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Error
	if err != nil {
		return &[]Like{}, err
	}
	return &likes, err
}

// DeleteUserLikes deletes all likes associated with the given user ID.
// It also updates the respective posts' Likes field and profile's LikedPosts field.
// DeleteUserLikes deletes all likes associated with the given user ID.
// It also updates the respective posts' Likes field and profile's LikedPosts field.
// Finally, it deletes the Like model entries for the given user.
func (l *Like) DeleteUserLikes(db *gorm.DB, profileID uint32) (int64, error) {
	// Find all the likes associated with the given user ID.
	var likes []*Like
	if err := db.Model(&Like{}).Where("profile_id = ?", profileID).Find(&likes).Error; err != nil {
		return 0, err
	}

	// Prepare a list of post IDs that will need to be updated.
	postIDs := make([]uint, 0)

	// Delete the likes and update the post IDs list.
	for _, like := range likes {
		postIDs = append(postIDs, like.PostID)
		if err := db.Delete(like).Error; err != nil {
			return 0, err
		}
	}

	// Update the LikedPosts field of the user's profile to remove the posts with deleted likes.
	// profile := Profile{}
	// if err := db.Preload("LikedPosts").Where("profile_id = ?", profileID).Take(&profile).Error; err != nil {
	// 	return 0, err
	// }

	// // Filter out the posts that have been deleted.
	// filteredLikedPosts := make([]*Post, 0, len(profile.LikedPosts))
	// for _, post := range profile.LikedPosts {
	// 	if !containsUint32(postIDs, uint32(post.ID)) {
	// 		filteredLikedPosts = append(filteredLikedPosts, post)
	// 	}
	// }
	// profile.LikedPosts = filteredLikedPosts

	// // Save the updated profile.
	// if err := db.Save(&profile).Error; err != nil {
	// 	return 0, err
	// }

	// Delete the Like model entries for the given user.
	if err := db.Where("profile_id = ?", profileID).Delete(&Like{}).Error; err != nil {
		return 0, err
	}

	return db.RowsAffected, nil
}

// DeletePostLikes deletes all likes associated with the given post ID.
// It also updates the respective users' LikedPosts field to remove the deleted post.
func (l *Like) DeletePostLikes(db *gorm.DB, pid uint64) (int64, error) {
	// Find all the likes associated with the given post ID.
	var likes []*Like
	if err := db.Model(&Like{}).Where("post_id = ?", pid).Find(&likes).Error; err != nil {
		return 0, err
	}

	// Prepare a list of user IDs that will need to be updated.
	userIDs := make([]uint, 0)

	// Delete the likes and update the user IDs list.
	for _, like := range likes {
		userIDs = append(userIDs, like.ProfileID)
		if err := db.Delete(like).Error; err != nil {
			return 0, err
		}
	}

	// Delete the Like model entries for the given post.
	if err := db.Where("post_id = ?", pid).Delete(&Like{}).Error; err != nil {
		return 0, err
	}

	return db.RowsAffected, nil
}
