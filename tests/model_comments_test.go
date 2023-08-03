package tests

import (
	"log"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	profile, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatalf("Cannot seed user and post %v\n", err)
	}

	newComment := models.Comment{
		Body:      "This is the comment body",
		ProfileID: uint32(profile.ID),
		PostID:    uint64(post.ID),
	}

	savedComment, err := newComment.SaveComment(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the comment: %v\n", err)
		return
	}

	assert.Equal(t, newComment.ProfileID, savedComment.ProfileID)
	assert.Equal(t, newComment.PostID, savedComment.PostID)
	assert.Equal(t, newComment.Body, "This is the comment body")
}

func TestCommentsForAPost(t *testing.T) {
	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	post, profiles, comments, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}
	//Where commentInstance is an instance of the post initialize in setup_test.go
	_, err = commentInstance.GetComments(server.DB, uint64(post.ID))
	if err != nil {
		t.Errorf("this is the error getting the comments: %v\n", err)
		return
	}
	assert.Equal(t, len(comments), 2)
	assert.Equal(t, len(profiles), 2)
}

func TestDeleteAComment(t *testing.T) {
	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	_, _, comments, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}
	// Delete the first comment
	for _, v := range comments {
		if v.ID == 2 {
			continue
		}
		commentInstance.ID = v.ID //commentInstance is defined in setup_test.go
	}
	isDeleted, err := commentInstance.DeleteAComment(server.DB)
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteCommentsForAPost(t *testing.T) {
	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	post, _, _, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}
	numberDeleted, err := commentInstance.DeleteUserComments(server.DB, uint32(post.ID))
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}

func TestDeleteCommentsForAUser(t *testing.T) {
	var profileID uint32

	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post and comment table %v\n", err)
	}
	_, profiles, _, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}

	// get the first user. When you delete this user, also delete his comment
	for _, v := range profiles {
		if v.ID == 2 {
			continue
		}
		profileID = uint32(v.ID)
	}
	numberDeleted, err := commentInstance.DeleteUserComments(server.DB, profileID)
	if err != nil {
		t.Errorf("this is the error deleting the comment: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}
