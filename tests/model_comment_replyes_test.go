package tests

import (
	"log"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateCommentReplye(t *testing.T) {
	err := refreshUserProfilePostAndCommentReplyeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post, comment, and replye table %v\n", err)
	}

	post, _, comments, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment table %v\n", err)
	}

	profile, err := seedOneUserProfile()

	if err != nil {
		log.Fatalf("cannot seed user profile table: %v", err)
	}

	// Get the first comment
	for _, v := range comments {
		if v.ID == 2 {
			continue
		}
		commentInstance.ID = v.ID //commentInstance is defined in setup_test.go
	}

	newCommentReplye := models.Replyes{
		Body:      "This is the comment replye body",
		CommentID: uint64(commentInstance.ID),
		ProfileID: uint64(profile.ID),
		PostID:    uint32(post.ID),
	}

	savedCommentReplye, err := newCommentReplye.SaveCommentReplyes(server.DB)
	if err != nil {
		t.Errorf("this is the error saved the comment replye: %v\n", err)
		return
	}

	assert.Equal(t, newCommentReplye.ProfileID, savedCommentReplye.ProfileID)
	assert.Equal(t, newCommentReplye.CommentID, savedCommentReplye.CommentID)
	assert.Equal(t, newCommentReplye.PostID, savedCommentReplye.PostID)
	assert.Equal(t, newCommentReplye.Body, savedCommentReplye.Body)
}

func TestCommentReplyeForAPost(t *testing.T) {
	err := refreshUserProfilePostAndCommentReplyeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post, comment, and replye table %v\n", err)
	}

	post, profiles, comment, replyes, err := seedUsersProfilePostsAndCommentReplyes()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment replye table %v\n", err)
	}

	//Where commentInstance is an instance of the post initialize in setup_test.go
	_, err = commentReplyesInstance.GetCommentReplyes(server.DB, uint64(comment.ID))
	if err != nil {
		t.Errorf("this is the error getting the comments: %v\n", err)
		return
	}

	assert.Equal(t, len(replyes), 2)
	assert.Equal(t, len(profiles), 2)

	for i, v := range replyes {
		assert.Equal(t, v.PostID, uint32(post.ID))
		assert.Equal(t, v.CommentID, uint64(comment.ID))
		assert.Equal(t, v.ProfileID, uint64(profiles[i].ID))
	}
}

func TestDeleteACommentReplye(t *testing.T) {
	err := refreshUserProfilePostAndCommentReplyeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post, comment, and replye table %v\n", err)
	}

	_, _, _, replyes, err := seedUsersProfilePostsAndCommentReplyes()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment replye table %v\n", err)
	}

	// Delete the first comment
	for _, v := range replyes {
		if v.ID == 2 {
			continue
		}
		commentReplyesInstance.ID = v.ID //commentInstance is defined in setup_test.go
	}

	isDeleted, err := commentReplyesInstance.DeleteAReplyes(server.DB)
	if err != nil {
		t.Errorf("this is the error deleting the comment: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteCommentReplyesForAPost(t *testing.T) {
	err := refreshUserProfilePostAndCommentReplyeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post, comment, and replye table %v\n", err)
	}

	post, _, _, _, err := seedUsersProfilePostsAndCommentReplyes()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment replye table %v\n", err)
	}

	numberDeleted, err := commentReplyesInstance.DeletePostCommentReplyes(server.DB, uint64(post.ID))
	if err != nil {
		t.Errorf("this is the error deleting the like: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(2))
}

func TestDeleteCommentReplyeForAUser(t *testing.T) {
	var profileID uint32
	err := refreshUserProfilePostAndCommentReplyeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, post, comment, and replye table %v\n", err)
	}

	_, profiles, _, _, err := seedUsersProfilePostsAndCommentReplyes()
	if err != nil {
		log.Fatalf("Error seeding user, post and comment replye table %v\n", err)
	}

	// get the first user. When you delete this user, also delete his comment
	for _, v := range profiles {
		if v.ID == 2 {
			continue
		}
		profileID = uint32(v.ID)
	}

	numberDeleted, err := commentReplyesInstance.DeleteUserProfileCommentReplyes(server.DB, profileID)
	if err != nil {
		t.Errorf("this is the error deleting the comment: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}
