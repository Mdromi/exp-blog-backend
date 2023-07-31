package tests

import (
	"log"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveALike(t *testing.T) {
	err := refreshUserProfilePostAndLikeTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile, post and like table %v\n", err)
	}

	profile, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatalf("Cannot seed profile and post %v\n", err)
	}

	newLike := models.Like{
		ProfileID: profile.ID,
		PostID:    post.ID,
	}

	savedLike, err := newLike.SaveLike(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the like: %v\n", err)
		return
	}
	assert.Equal(t, newLike.ID, savedLike.ID)
	assert.Equal(t, newLike.ProfileID, savedLike.ProfileID)
	assert.Equal(t, newLike.PostID, savedLike.PostID)

	// Refresh database all table
	err = refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
}
