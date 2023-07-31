package models

import (
	"errors"

	"gorm.io/gorm"
)

// Dislike model represents a dislike by a user on a post
type Dislike struct {
	gorm.Model
	ProfileID uint `gorm:"not null" json:"profile_id"`
	PostID    uint `gorm:"not null" json:"post_id"`
}

// TASK: uid to convert profileID | likes & dislikes

func (l *Like) SaveDislike(db *gorm.DB) (*Like, error) {
	// Check if the user has previously liked this post
	likedBefore, err := CheckIfLikedBefore(db, l.PostID, l.ProfileID)
	if err != nil {
		return &Like{}, err
	}

	if likedBefore {
		// The user has previously liked this post, so we need to remove the like
		err = RemoveLike(db, l.PostID, l.ProfileID)
		if err != nil {
			return &Like{}, err
		}
	}

	// Check if the user has disliked this post before
	dislikedBefore, err := CheckIfDislikedBefore(db, l.PostID, l.ProfileID)
	if err != nil {
		return &Like{}, err
	}

	if !dislikedBefore {
		// The user has not disliked this post before, so let's save the incoming dislike
		err = SaveNewDislike(db, l.PostID, l.ProfileID)
		if err != nil {
			return &Like{}, err
		}

	} else {
		// The user has disliked it before, so return a custom error message
		return &Like{}, errors.New("double dislike")
	}

	return l, nil
}

func (d *Dislike) DeleteDislike(db *gorm.DB) (*Dislike, error) {
	// Check if the user has disliked this post before
	dislikedBefore, err := CheckIfDislikedBefore(db, d.PostID, d.ProfileID)
	if err != nil {
		return &Dislike{}, err
	}

	if dislikedBefore {
		// The user has disliked this post, so we need to remove the dislike
		err = RemoveDislike(db, d.PostID, d.ProfileID)
		if err != nil {
			return &Dislike{}, err
		}

		// Update the profile's DislikedPosts field to remove the post
		// err = RemoveFromDislikedPosts(db, d.PostID, d.ProfileID)
		// if err != nil {
		// 	return &Dislike{}, err
		// }
	} else {
		// The user has not disliked this post before, so return a custom error message
		return &Dislike{}, errors.New("dislike not found")
	}

	return d, nil
}

func (l *Like) GetDislikesInfo(db *gorm.DB, pid uint64) (*[]Dislike, error) {
	dislike := []Dislike{}

	err := db.Debug().Model(&Like{}).Where("post_id = ?", pid).Find(&dislike).Error
	if err != nil {
		return &[]Dislike{}, err
	}
	return &dislike, err
}

// DeleteUserDislikes deletes all dislikes associated with the given user ID.
// It also updates the respective posts' Dislikes field and profile's DislikedPosts field.
// Finally, it deletes the Dislike model entries for the given user.
func (d *Dislike) DeleteUserDislikes(db *gorm.DB, profileID uint) (int64, error) {
	// Find all the dislikes associated with the given user ID.
	var dislikes []*Dislike
	if err := db.Model(&Dislike{}).Where("profile_id = ?", profileID).Find(&dislikes).Error; err != nil {
		return 0, err
	}

	// Prepare a list of post IDs that will need to be updated.
	postIDs := make([]uint, 0)

	// Delete the dislikes and update the post IDs list.
	for _, dislike := range dislikes {
		postIDs = append(postIDs, dislike.PostID)
		if err := db.Delete(dislike).Error; err != nil {
			return 0, err
		}
	}

	// Update the DislikedPosts field of the user's profile to remove the posts with deleted dislikes.
	// profile := Profile{}
	// if err := db.Preload("DislikedPosts").Where("profile_id = ?", profileID).Take(&profile).Error; err != nil {
	// 	return 0, err
	// }

	// // Filter out the posts that have been deleted.
	// filteredDislikedPosts := make([]*Post, 0, len(profile.DislikedPosts))
	// for _, post := range profile.DislikedPosts {
	// 	if !containsUint32(postIDs, uint32(post.ID)) {
	// 		filteredDislikedPosts = append(filteredDislikedPosts, post)
	// 	}
	// }
	// profile.DislikedPosts = filteredDislikedPosts

	// // Save the updated profile.
	// if err := db.Save(&profile).Error; err != nil {
	// 	return 0, err
	// }

	// Delete the Dislike model entries for the given user.
	if err := db.Where("profile_id = ?", profileID).Delete(&Dislike{}).Error; err != nil {
		return 0, err
	}

	return db.RowsAffected, nil
}

// DeletePostDislikes deletes all dislikes associated with the given post ID.
// It also updates the respective users' DislikedPosts field to remove the deleted post.
func (d *Dislike) DeletePostDislikes(db *gorm.DB, pid uint64) (int64, error) {
	// Find all the dislikes associated with the given post ID.
	var dislikes []*Dislike
	if err := db.Model(&Dislike{}).Where("post_id = ?", pid).Find(&dislikes).Error; err != nil {
		return 0, err
	}

	// Prepare a list of user IDs that will need to be updated.
	userIDs := make([]uint, 0)

	// Delete the dislikes and update the user IDs list.
	for _, dislike := range dislikes {
		userIDs = append(userIDs, dislike.ProfileID)
		if err := db.Delete(dislike).Error; err != nil {
			return 0, err
		}
	}

	// Update the DislikedPosts field of respective users to remove the post with deleted dislikes.
	// for _, userID := range userIDs {
	// 	if err := RemoveFromDislikedPosts(db, uint(pid), userID); err != nil {
	// 		return 0, err
	// 	}
	// }

	// Delete the Dislike model entries for the given post.
	if err := db.Where("post_id = ?", pid).Delete(&Dislike{}).Error; err != nil {
		return 0, err
	}

	return db.RowsAffected, nil
}
