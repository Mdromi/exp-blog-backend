package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/stretchr/testify/assert"
)

func TestFindAllPosts(t *testing.T) {
	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile and post table %v\n", err)
	}

	_, _, err = seedUsersProfileAndPosts()
	if err != nil {
		log.Fatalf("Error seeding user and post table %v\n", err)
	}

	// Where postInstance is an instance of the post initialize in setup_test.go
	posts, err := postInstance.FindAllPosts(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the posts: %v\n", err)
		return
	}
	assert.Equal(t, len(*posts), 2)

	// Refresh database all table
	err = refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
}

func TestSavePost(t *testing.T) {
	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile and post table %v\n", err)
	}

	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatalf("Cannot seed user profile %v\n", err)
	}

	newPost := models.Post{
		Title:    "This is the title",
		Content:  "This is the content",
		AuthorID: profile.ID,
		Author:   profile,
	}

	savedPost, err := newPost.SavePost(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the post: %v\n", err)
		return
	}
	fmt.Println("newPost.ID", newPost.ID)
	fmt.Println("savedPost.ID", savedPost.ReadTime)

	assert.Equal(t, newPost.ID, savedPost.ID)
	assert.Equal(t, newPost.Title, savedPost.Title)
	assert.Equal(t, newPost.Content, savedPost.Content)
	assert.Equal(t, newPost.AuthorID, savedPost.AuthorID)

	// Refresh database all table
	err = refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
}

func TestFindPostByID(t *testing.T) {
	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile and post table %v\n", err)
	}

	_, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	foundPost, err := post.FindPostById(server.DB, uint64(post.ID))
	if err != nil {
		t.Errorf("this is the error getting one user: %v\n", err)
		return
	}
	assert.Equal(t, foundPost.ID, post.ID)
	assert.Equal(t, foundPost.Title, post.Title)
	assert.Equal(t, foundPost.Content, post.Content)

	// Refresh database all table
	err = refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
}

func TestUpdateAPost(t *testing.T) {
	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile and post table %v\n", err)
	}

	_, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	postUpdate := models.Post{
		Title:    "modiUpdate",
		Content:  "modiupdate@example.com",
		AuthorID: post.AuthorID,
	}
	updatedPost, err := postUpdate.UpdateAPost(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, updatedPost.ID, postUpdate.ID)
	assert.Equal(t, updatedPost.Title, postUpdate.Title)
	assert.Equal(t, updatedPost.Content, postUpdate.Content)
	assert.Equal(t, updatedPost.AuthorID, postUpdate.AuthorID)

	// Refresh database all table
	err = refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
}

func TestDeleteAPost(t *testing.T) {
	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile and post table %v\n", err)
	}

	_, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	isDeleted, err := post.DeleteAPost(server.DB)
	if err != nil {
		t.Errorf("this is the error updating the user: %v\n", err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}

func TestDeleteUserPosts(t *testing.T) {
	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatalf("Error refreshing user, profile and post table %v\n", err)
	}

	_, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatalf("Error Seeding table")
	}

	numberDeleted, err := postInstance.DeleteUserPosts(server.DB, uint32(post.AuthorID))
	if err != nil {
		t.Errorf("this is the error deleting the post: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))

	// Refresh database all table
	err = refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
}
