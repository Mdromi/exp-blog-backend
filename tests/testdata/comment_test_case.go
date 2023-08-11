package testdata

import (
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/models"
)

type CreateCommentTestCase struct {
	PostIDString string
	InputJSON    string
	StatusCode   int
	ProfileID    uint32
	PostID       uint
	Body         string
	TokenGiven   string
}

type GetCommentTestCase struct {
	PostID         string
	ProfileLength  int
	CommentsLength int
	StatusCode     int
}

type UpdateCommentsTestCase struct {
	CommentID  string
	UpdateJSON string
	Body       string
	TokenGiven string
	StatusCode int
}

type DeleteCommentsTestCase struct {
	CommentID      string
	UsersLength    int
	TokenGiven     string
	CommentsLength int
	StatusCode     int
}

func CreateCommentsSamples(firstUserToken, secondUserToken string, firstPostID uint) []CreateCommentTestCase {
	return []CreateCommentTestCase{
		{
			// User 1 can comment on his post
			PostIDString: strconv.Itoa(int(firstPostID)), //we need the id as a string
			InputJSON:    `{"body": "comment from user 1"}`,
			StatusCode:   201,
			ProfileID:    1,
			PostID:       firstPostID,
			Body:         "comment from user 1",
			TokenGiven:   firstUserToken,
		},
		{
			// User 2 can also comment on user 1 post
			PostIDString: strconv.Itoa(int(firstPostID)),
			InputJSON:    `{"body":"comment from user 2"}`,
			StatusCode:   201,
			ProfileID:    2,
			PostID:       firstPostID,
			Body:         "comment from user 2",
			TokenGiven:   secondUserToken,
		},
		{
			// When no body is provided:
			PostIDString: strconv.Itoa(int(firstPostID)),
			InputJSON:    `{"body":""}`,
			StatusCode:   422,
			PostID:       firstPostID,
			TokenGiven:   secondUserToken,
		},
		{
			// Not authenticated (No token provided)
			PostIDString: strconv.Itoa(int(firstPostID)),
			StatusCode:   401,
			TokenGiven:   "",
		},
		{
			// Wrong Token
			PostIDString: strconv.Itoa(int(firstPostID)),
			StatusCode:   401,
			TokenGiven:   "This is an incorrect token",
		},
		{
			// When invalid post id is given
			PostIDString: "unknwon",
			StatusCode:   400,
		},
	}
}

func GetCommentsSamples(profiles []*models.Profile, comments []models.Comment, postID string) []GetCommentTestCase {
	return []GetCommentTestCase{
		{
			PostID:         postID,
			StatusCode:     200,
			ProfileLength:  len(profiles),
			CommentsLength: len(comments),
		},
		{
			PostID:     "unknwon",
			StatusCode: 400,
		},
		{
			PostID:     strconv.Itoa(12322), //an id that does not exist
			StatusCode: 404,
		},
	}
}

func UpdateCommentsSamples(tokenString, secondCommentID string) []UpdateCommentsTestCase {
	return []UpdateCommentsTestCase{
		{
			CommentID:  secondCommentID,
			UpdateJSON: `{"Body":"This is the update body"}`,
			StatusCode: 200,
			Body:       "This is the update body",
			TokenGiven: tokenString,
		},
		{
			// When the body field is empty
			CommentID:  secondCommentID,
			UpdateJSON: `{"Body":""}`,
			StatusCode: 422,
			TokenGiven: tokenString,
		},
		{
			//an id that does not exist
			CommentID:  strconv.Itoa(12322),
			StatusCode: 404,
			TokenGiven: tokenString,
		},
		{
			//When the user is not authenticated
			CommentID:  secondCommentID,
			StatusCode: 401,
			TokenGiven: "",
		},
		{
			//When wrong token is passed
			CommentID:  secondCommentID,
			StatusCode: 401,
			TokenGiven: "this is a wrong token",
		},
		{
			// When id passed is invalid
			CommentID:  "unknwon",
			StatusCode: 400,
		},
	}
}

func DeleteCommentsSamples(tokenString, secondCommentID string) []DeleteCommentsTestCase {
	return []DeleteCommentsTestCase{
		{
			CommentID:  secondCommentID,
			StatusCode: 200,
			TokenGiven: tokenString,
		},
		{
			//an id that does not exist
			CommentID:  strconv.Itoa(12322),
			StatusCode: 404,
			TokenGiven: tokenString,
		},
		{
			//When the user is not authenticated
			CommentID:  secondCommentID,
			StatusCode: 401,
			TokenGiven: "",
		},
		{
			//When wrong token is passed
			CommentID:  secondCommentID,
			StatusCode: 401,
			TokenGiven: "this is a wrong token",
		},
		{
			// When id passed is invalid
			CommentID:  "unknwon",
			StatusCode: 400,
			TokenGiven: tokenString,
		},
	}
}
