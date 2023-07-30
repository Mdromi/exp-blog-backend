package models

import (
	"errors"

	"gorm.io/gorm"
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

// RemoveFromDislikedPosts removes a post from the user's DislikedPosts list.
func RemoveFromDislikedPosts(db *gorm.DB, postID, profileID uint) error {
	profile := Profile{}
	err := db.Debug().Preload("DislikedPosts").Where("id = ?", profileID).Take(&profile).Error
	if err != nil {
		return err
	}

	// Find the index of the disliked post in DislikedPosts
	index := -1
	for i, post := range profile.DislikedPosts {
		if post.ID == postID {
			index = i
			break
		}
	}

	// Remove the post from DislikedPosts if found
	if index != -1 {
		profile.DislikedPosts = append(profile.DislikedPosts[:index], profile.DislikedPosts[index+1:]...)
		return db.Debug().Save(&profile).Error
	}

	return nil
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

// AddToLikedPosts adds a post to the user's LikedPosts list.
func AddToLikedPosts(db *gorm.DB, postID, profileID uint) error {
	profile := Profile{}
	err := db.Debug().Preload("LikedPosts").Where("id = ?", profileID).Take(&profile).Error
	if err != nil {
		return err
	}

	// Check if the post already exists in the LikedPosts
	for _, post := range profile.LikedPosts {
		if post.ID == postID {
			return errors.New("double like")
		}
	}

	// Append the new liked post to the LikedPosts field
	profile.LikedPosts = append(profile.LikedPosts, &Post{Model: gorm.Model{ID: postID}})
	return db.Debug().Save(&profile).Error
}

// RemoveLike removes a like entry from the database for a post by a user.
func RemoveLike(db *gorm.DB, postID, profileID uint) error {
	like := Like{}
	return db.Debug().Delete(&like, "post_id = ? AND profile_id = ?", postID, profileID).Error
}

// RemoveFromLikedPosts removes a post from the user's LikedPosts list.
func RemoveFromLikedPosts(db *gorm.DB, postID, profileID uint) error {
	profile := Profile{}
	err := db.Debug().Preload("LikedPosts").Where("id = ?", profileID).Take(&profile).Error
	if err != nil {
		return err
	}

	// Find the index of the liked post in LikedPosts
	index := -1
	for i, post := range profile.LikedPosts {
		if post.ID == postID {
			index = i
			break
		}
	}

	// Remove the post from LikedPosts if found
	if index != -1 {
		profile.LikedPosts = append(profile.LikedPosts[:index], profile.LikedPosts[index+1:]...)
		return db.Debug().Save(&profile).Error
	}

	return nil
}

// SaveNewDislike saves a new dislike entry to the database for a post by a user.
func SaveNewDislike(db *gorm.DB, postID, profileID uint) error {
	dislike := &Dislike{
		ProfileID: profileID,
		PostID:    postID,
	}
	return db.Debug().Create(dislike).Error
}

// AddToDislikedPosts adds a post to the user's DislikedPosts list.
func AddToDislikedPosts(db *gorm.DB, postID, profileID uint) error {
	profile := Profile{}
	err := db.Debug().Preload("DislikedPosts").Where("id = ?", profileID).Take(&profile).Error
	if err != nil {
		return err
	}

	// Check if the post already exists in the DislikedPosts
	for _, post := range profile.DislikedPosts {
		if post.ID == postID {
			return errors.New("double dislike")
		}
	}

	// Append the new disliked post to the DislikedPosts field
	profile.DislikedPosts = append(profile.DislikedPosts, &Post{Model: gorm.Model{ID: postID}})
	return db.Debug().Save(&profile).Error
}

// UpdatePostLikeDislike updates the post's Likes and Dislikes fields based on the 'like' boolean.
func UpdatePostLikeDislike(db *gorm.DB, postID, profileID uint, like bool) error {
	post := Post{}
	err := db.Debug().Preload("Likes").Preload("Dislikes").Where("id = ?", postID).Take(&post).Error
	if err != nil {
		return err
	}

	if like {
		// Check if the user has already liked the post
		for _, l := range post.Likes {
			if l.ProfileID == profileID {
				return errors.New("double like")
			}
		}

		// Create a new Like model and append it to the Likes field
		newLike := &Like{ProfileID: profileID, PostID: postID}
		post.Likes = append(post.Likes, newLike)
	} else {
		// Check if the user has already disliked the post
		for _, dl := range post.Dislikes {
			if dl.ProfileID == profileID {
				return errors.New("double dislike")
			}
		}

		// Create a new Dislike model and append it to the Dislikes field
		newDislike := &Dislike{ProfileID: profileID, PostID: postID}
		post.Dislikes = append(post.Dislikes, newDislike)
	}

	// Save the updated post
	return db.Debug().Save(&post).Error
}

// RemoveFromPostLikes removes a like from the Likes field of the Post model.
func RemoveFromPostLikes(db *gorm.DB, postID, profileID uint) error {
	post := Post{}
	err := db.Debug().Preload("Likes").Where("id = ?", postID).Take(&post).Error
	if err != nil {
		return err
	}

	// Find the index of the like in Likes
	index := -1
	for i, like := range post.Likes {
		if like.ProfileID == profileID {
			index = i
			break
		}
	}

	// Remove the like from Likes if found
	if index != -1 {
		post.Likes = append(post.Likes[:index], post.Likes[index+1:]...)
		return db.Debug().Save(&post).Error
	}

	return nil
}

// AddToPostDislikes adds a dislike to the Dislikes field of the Post model.
func AddToPostDislikes(db *gorm.DB, postID, profileID uint) error {
	post := Post{}
	err := db.Debug().Preload("Dislikes").Where("id = ?", postID).Take(&post).Error
	if err != nil {
		return err
	}

	// Check if the dislike already exists in the Dislikes
	for _, dl := range post.Dislikes {
		if dl.ProfileID == profileID {
			return errors.New("double dislike")
		}
	}

	// Create a new Dislike model and append it to the Dislikes field
	newDislike := &Dislike{ProfileID: profileID, PostID: postID}
	post.Dislikes = append(post.Dislikes, newDislike)

	return db.Debug().Save(&post).Error
}

// RemoveFromPostDislikes removes the dislike from the Dislikes field of the Post model.
func RemoveFromPostDislikes(db *gorm.DB, postID, dislikeID uint) error {
	post := Post{}
	err := db.Debug().Preload("Dislikes").Where("id = ?", postID).Take(&post).Error
	if err != nil {
		return err
	}

	// Find the index of the dislike in Dislikes
	index := -1
	for i, dislike := range post.Dislikes {
		if dislike.ID == dislikeID {
			index = i
			break
		}
	}

	// Remove the dislike from Dislikes if found
	if index != -1 {
		post.Dislikes = append(post.Dislikes[:index], post.Dislikes[index+1:]...)
		return db.Debug().Save(&post).Error
	}

	return nil
}
