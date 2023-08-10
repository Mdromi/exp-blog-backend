package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	executeablefunctions "github.com/Mdromi/exp-blog-backend/tests/executeable_functions"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var ExecuteCreatePostTest = executeablefunctions.ExecuteCreatePostTest
var ExecuteGetPostByID = executeablefunctions.ExecuteGetPostByID
var ExecuteUpdatePost = executeablefunctions.ExecuteUpdatePost
var ExecuteDeletePost = executeablefunctions.ExecuteDeletePost

func TestCreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, tokenString := seedProfileAndSignIn(server.DB)

	// Get test samples for creating post and iterate over them.
	samples := testdata.CreatePostsSamples(tokenString)
	ExecuteCreatePostTest(t, samples, &server)
}

func TestGetPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = seedUsersProfileAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/posts", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	postsInterface := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.String()), &postsInterface)
	// This is so that we can get the length of the posts:
	thePosts := postsInterface["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(thePosts), 2)
}

func TestGetPostByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	_, post, err := seedOneUserProfileAndOnePost()
	if err != nil {
		log.Fatal(err)
	}

	// Get test samples for creating post and iterate over them.
	samples := testdata.GetPostByIDSamples(post)
	ExecuteGetPostByID(t, samples, &server)
}

func TestUpdatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Author token and authPostID
	tokenString, AuthPostID := getTokenAndPostID()

	// Get test samples for updating post and iterate over them.
	samples := testdata.UpdatePostTestSamples(tokenString, AuthPostID)
	ExecuteUpdatePost(t, samples, &server)
}

func TestDeletePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Author token and authPostID
	tokenString, AuthPostID := getTokenAndPostID()

	// Get test samples for updating post and iterate over them.
	samples := testdata.DeletePostTestSamples(tokenString, AuthPostID)
	ExecuteDeletePost(t, samples, &server)
}

func getTokenAndPostID() (string, uint) {
	var PostUserEmail, PostUserPassword string
	var AuthPostID uint

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	profiles, posts, err := seedUsersProfileAndPosts()
	if err != nil {
		log.Fatal(err)
	}

	// Get Only the first profile
	for _, profile := range profiles {
		if profile.ID == 2 {
			continue
		}
		user, err := controllers.FindUserByID(server.DB, uint32(profile.UserID))
		if err != nil {
			log.Fatal(err)
		}
		PostUserEmail = user.Email
		PostUserPassword = "password" // Note the password in the database is already hashed, we want unhashed
	}

	// Get only the first post
	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		AuthPostID = post.ID
	}

	tokenInterface, err := server.SignIn(PostUserEmail, PostUserPassword)
	if err != nil {
		log.Fatal(err)
	}

	token := tokenInterface["token"] // get only the token
	tokenString := fmt.Sprintf("Bearer %v", token.(string))

	return tokenString, AuthPostID
}
