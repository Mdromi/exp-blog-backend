package tests

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/api/models"
	executeablefunctions "github.com/Mdromi/exp-blog-backend/tests/executeable_functions"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
)

var ExecuteCreateComments = executeablefunctions.ExecuteCreateComments
var ExecuteGetComments = executeablefunctions.ExecuteGetComments
var ExecuteUpdateComments = executeablefunctions.ExecuteUpdateComments
var ExecuteDeleteComments = executeablefunctions.ExecuteDeleteComments

func TestCommnetPost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var firstPostID uint

	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}

	// Author token and authPostID
	firstUserToken, secondUserToken, firstPostID := getUserTokensAndPostIDForComments()

	// Get test samples for updating post and iterate over them.
	samples := testdata.CreateCommentsSamples(firstUserToken, secondUserToken, firstPostID)
	ExecuteCreateComments(t, samples, &server)
}

func TestGetComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}
	post, profiles, comments, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}

	postID := strconv.Itoa(int(post.ID))

	// Get test samples for updating post and iterate over them.
	samples := testdata.GetCommentsSamples(profiles, comments, postID)
	ExecuteGetComments(t, samples, &server)
}

func TestUpdateComment(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var secondUserEmail, password string

	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}
	post, profiles, comments, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}

	postID := float64(post.ID)

	// Get profile
	secondUserID, secondUserEmail, err := extractUserInfoForProfileID(&server, profiles, 2)
	if err != nil {
		log.Fatal(err)
	}

	password = "password"

	// Get only the second comment
	secondCommentID, err := extractCommentID(comments, 2) // Target comment with ID 2
	if err != nil {
		log.Fatal(err)
	}

	//Login the user and get the authentication token
	secondUserToken, err := getUserAuthToken(&server, secondUserEmail, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	// Get test samples for updating post and iterate over them.
	samples := testdata.UpdateCommentsSamples(secondUserToken, secondCommentID)
	ExecuteUpdateComments(t, samples, &server, postID, secondUserID)
}

func TestDeleteComment(t *testing.T) {

	gin.SetMode(gin.TestMode)

	var secondUserEmail, password, secondCommentID string

	err := refreshUserProfilePostAndCommentTable()
	if err != nil {
		log.Fatal(err)
	}

	post, profiles, comments, err := seedUsersProfilePostsAndComments()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}
	// Get only the second profile user
	_, secondUserEmail, err = extractUserInfoForProfileID(&server, profiles, 2)
	if err != nil {
		log.Fatal(err)
	}
	password = "password"
	postID := float64(post.ID)

	// Get only the second comment
	secondCommentID, err = extractCommentID(comments, 2) // Target comment with ID 2
	if err != nil {
		log.Fatal(err)
	}

	//Login the user and get the authentication token
	tokenString, err := getUserAuthToken(&server, secondUserEmail, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	// Get test samples for updating post and iterate over them.
	samples := testdata.DeleteCommentsSamples(tokenString, secondCommentID)
	ExecuteDeleteComments(t, samples, &server, postID)

}

func getUserTokensAndPostIDForComments() (string, string, uint) {
	var firstUserEmail, secondUserEmail string
	var firstPostID uint

	err := refreshUserProfileAndPostTable()
	if err != nil {
		log.Fatal(err)
	}

	profiles, posts, err := seedUsersProfileAndPosts()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	// Get profile
	for _, profile := range profiles {
		if profile.ID == 1 {
			user, err := controllers.FindUserByID(server.DB, uint32(profile.UserID))
			if err != nil {
				log.Fatal(err)
			}
			firstUserEmail = user.Email
		}

		if profile.ID == 2 {
			user, err := controllers.FindUserByID(server.DB, uint32(profile.UserID))
			if err != nil {
				log.Fatal(err)
			}

			secondUserEmail = user.Email
		}

	}

	// Get only the first post, which belongs to first user
	for _, post := range posts {
		if post.ID == 2 {
			continue
		}
		firstPostID = post.ID
	}

	// Login both users
	// user 1 and user 2 password are the same, you can change if you want (Note by the time they are hashed and saved in the db, they are different)
	// Note: the value of the user password before it was hashed is "password". so:
	password := "password"

	// Login First User
	tokenInterface1, err := server.SignIn(firstUserEmail, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token1 := tokenInterface1["token"] //get only the token
	firstUserToken := fmt.Sprintf("Bearer %v", token1)

	// Login Second User
	tokenInterface2, err := server.SignIn(secondUserEmail, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token2 := tokenInterface2["token"] //get only the token
	secondUserToken := fmt.Sprintf("Bearer %v", token2)

	return firstUserToken, secondUserToken, firstPostID
}

func getUserAuthToken(server *controllers.Server, userEmail, password string) (string, error) {
	tokenInterface, err := server.SignIn(userEmail, password)
	if err != nil {
		return "", err
	}
	token := tokenInterface["token"]
	return fmt.Sprintf("Bearer %v", token), nil
}

func extractUserInfoForProfileID(server *controllers.Server, profiles []*models.Profile, targetProfileID uint) (float64, string, error) {
	var userID float64
	var userEmail string

	for _, profile := range profiles {
		if profile.ID == targetProfileID {
			user, err := controllers.FindUserByID(server.DB, uint32(profile.UserID))
			if err != nil {
				return 0, "", err
			}

			userID = float64(user.ID)
			userEmail = user.Email
			break
		}
	}

	return userID, userEmail, nil
}

func extractCommentID(comments []models.Comment, targetID int) (string, error) {
	var targetCommentID string

	for _, comment := range comments {
		if int(comment.ID) == targetID {
			targetCommentID = strconv.Itoa(int(comment.ID))
			break
		}
	}

	if targetCommentID == "" {
		return "", fmt.Errorf("comment with ID %d not found", targetID)
	}

	return targetCommentID, nil
}
