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

var ExecuteCreateCommentReplyes = executeablefunctions.ExecuteCreateCommentReplyes
var ExecuteGetCommentReplyes = executeablefunctions.ExecuteGetCommentReplyes
var ExecuteUpdateCommentReplye = executeablefunctions.ExecuteUpdateCommentReplye
var ExecuteDeleteCommentReplye = executeablefunctions.ExecuteDeleteCommentReplye

func TestCommnetReplyePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	// Author token and authPostID
	firstUserToken, secondUserToken, secondCommentID, _, firstPostID := getUserTokensAndPostIDForCommentReplyes()

	err = refreshReplyeTable()
	if err != nil {
		log.Fatal(err)
	}
	// Get test samples for updating post and iterate over them.
	samples := testdata.CreateCommentReplyeSamples(firstUserToken, secondUserToken, secondCommentID, uint(firstPostID))
	ExecuteCreateCommentReplyes(t, samples, &server)
}

func TestGetCommentReplye(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}
	post, profiles, comment, replyes, err := seedUsersProfilePostsAndCommentReplyes()
	if err != nil {
		log.Fatalf("Cannot seed tables %v\n", err)
	}

	commentID := strconv.Itoa(int(comment.ID))
	postID := strconv.Itoa(int(post.ID))

	// Get test samples for updating post and iterate over them.
	samples := testdata.GetCommentReplyeSamples(profiles, replyes, postID, commentID)
	ExecuteGetCommentReplyes(t, samples, &server)
}

func TestUpdateCommentReplye(t *testing.T) {

	gin.SetMode(gin.TestMode)

	// Author token and authPostID
	_, secondUserToken, secondCommentID, secondReplyesID, firstPostID := getUserTokensAndPostIDForCommentReplyes()

	postID := strconv.Itoa(int(firstPostID))

	// Get test samples for updating post and iterate over them.
	samples := testdata.UpdateCommentReplyeSamples(secondUserToken, secondCommentID, postID, secondReplyesID)
	ExecuteUpdateCommentReplye(t, samples, &server)
}

func TestDeleteCommentReplye(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Author token and authPostID
	_, secondUserToken, secondCommentID, secondReplyesID, firstPostID := getUserTokensAndPostIDForCommentReplyes()

	postID := strconv.Itoa(int(firstPostID))
	// Get test samples for updating post and iterate over them.
	samples := testdata.DeleteCommentReplyeSamples(secondUserToken, secondCommentID, postID, secondReplyesID)
	ExecuteDeleteCommentReplye(t, samples, &server)
}

func getUserTokensAndPostIDForCommentReplyes() (string, string, string, string, uint) {
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	var firstUserEmail, secondUserEmail, secondCommentID, secondReplyesID string
	var firstPostID uint

	post, profiles, comment, replyes, err := seedUsersProfilePostsAndCommentReplyes()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}

	// Get profile
	emailMap, err := extractEmailsFromProfiles(profiles, &server)
	if err != nil {
		log.Fatal(err)
	}

	firstUserEmail = emailMap[0]
	secondUserEmail = emailMap[1]
	secondCommentID = strconv.Itoa(int(comment.ID))

	// Get only the first post, which belongs to first user
	for _, replye := range replyes {
		if replye.ID == 1 {
			continue
		}
		secondReplyesID = strconv.Itoa(int(replye.ID))
	}

	firstPostID = post.ID
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

	return firstUserToken, secondUserToken, secondCommentID, secondReplyesID, firstPostID
}

func extractEmailsFromProfiles(profiles []*models.Profile, server *controllers.Server) (map[int]string, error) {
	emailMap := make(map[int]string)

	for i, profile := range profiles {
		user, err := controllers.FindUserByID(server.DB, uint32(profile.UserID))
		if err != nil {
			return nil, err
		}
		emailMap[i] = user.Email
	}

	return emailMap, nil
}
