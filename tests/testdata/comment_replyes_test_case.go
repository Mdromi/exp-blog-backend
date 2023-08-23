package testdata

import (
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/models"
)

type CreateCommentReplyesTestCase struct {
	CommentID    string
	PostIDString string
	InputJSON    string
	StatusCode   int
	ProfileID    uint32
	PostID       uint
	Body         string
	TokenGiven   string
}

type GetCommentReplyeTestCase struct {
	CommentID     string
	PostID        string
	ProfileLength int
	ReplyesLength int
	StatusCode    int
}

type UpdateCommentReplyesTestCase struct {
	CommentID  string
	PostID     string
	ReplyesID  string
	UpdateJSON string
	Body       string
	TokenGiven string
	StatusCode int
}

type DeleteCommentReplyesTestCase struct {
	CommentID  string
	PostID     string
	ReplyesID  string
	TokenGiven string
	StatusCode int
}

func CreateCommentReplyeSamples(firstUserToken, secondUserToken, secondCommentID string, firstPostID uint) []CreateCommentReplyesTestCase {
	return []CreateCommentReplyesTestCase{
		{
			// User 1 can comment reply on his post
			CommentID:    secondCommentID,
			PostIDString: strconv.Itoa(int(firstPostID)), //we need the id as a string
			InputJSON:    `{"body": "comment replyes from user 1"}`,
			StatusCode:   201,
			ProfileID:    3,
			PostID:       firstPostID,
			Body:         "comment replyes from user 1",
			TokenGiven:   firstUserToken,
		},
		{
			// User 2 can also comment replye on user 1 post
			CommentID:    secondCommentID,
			PostIDString: strconv.Itoa(int(firstPostID)),
			InputJSON:    `{"body":"comment replyes from user 2"}`,
			StatusCode:   201,
			ProfileID:    4,
			PostID:       firstPostID,
			Body:         "comment replyes from user 2",
			TokenGiven:   secondUserToken,
		},
		{
			// When no body is provided:
			CommentID:    secondCommentID,
			PostIDString: strconv.Itoa(int(firstPostID)),
			InputJSON:    `{"body":""}`,
			StatusCode:   422,
			PostID:       firstPostID,
			TokenGiven:   secondUserToken,
		},
		{
			// Not authenticated (No token provided)
			CommentID:    secondCommentID,
			PostIDString: strconv.Itoa(int(firstPostID)),
			StatusCode:   401,
			TokenGiven:   "",
		},
		{
			// Wrong Token
			CommentID:    secondCommentID,
			PostIDString: strconv.Itoa(int(firstPostID)),
			StatusCode:   401,
			TokenGiven:   "This is an incorrect token",
		},
		{
			// When invalid post id is given
			CommentID:    secondCommentID,
			PostIDString: "unknwon",
			StatusCode:   401,
			TokenGiven:   secondUserToken,
		},
		{
			// When invalid comment id is given
			CommentID:    "unknwon",
			PostIDString: strconv.Itoa(int(firstPostID)),
			StatusCode:   400,
			TokenGiven:   secondUserToken,
		},
	}
}

func GetCommentReplyeSamples(profiles []*models.Profile, replyes []models.Replyes, postID, commentID string) []GetCommentReplyeTestCase {
	return []GetCommentReplyeTestCase{
		{
			CommentID:     commentID,
			PostID:        postID,
			StatusCode:    200,
			ProfileLength: len(profiles),
			ReplyesLength: len(replyes),
		},
		{
			CommentID:  "unknwon",
			PostID:     postID,
			StatusCode: 400,
		},
		{
			CommentID:  strconv.Itoa(12322), //an id that does not exist
			PostID:     postID,
			StatusCode: 404,
		},
	}
}

func UpdateCommentReplyeSamples(tokenString, secondCommentID, postID, replyesID string) []UpdateCommentReplyesTestCase {
	return []UpdateCommentReplyesTestCase{
		{
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			UpdateJSON: `{"Body":"This is the update body"}`,
			StatusCode: 200,
			Body:       "This is the update body",
			TokenGiven: tokenString,
		},
		{
			// When the body field is empty
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			UpdateJSON: `{"Body":""}`,
			StatusCode: 422,
			TokenGiven: tokenString,
		},
		{
			//an id that CommentID does not exist
			CommentID:  strconv.Itoa(12322),
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 404,
			TokenGiven: tokenString,
		},
		{
			//an id that ReplyesID does not exist
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  strconv.Itoa(12322),
			StatusCode: 400,
			TokenGiven: tokenString,
		},
		{
			//When the user is not authenticated
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 401,
			TokenGiven: "",
		},
		{
			//When wrong token is passed
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 401,
			TokenGiven: "this is a wrong token",
		},
		{
			// When id passed is invalid CommentID
			CommentID:  "unknwon",
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 400,
			TokenGiven: tokenString,
		},
		{
			// When id passed is invalid ReplyesID
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  "unknwon",
			StatusCode: 400,
			TokenGiven: tokenString,
		},
	}
}

func DeleteCommentReplyeSamples(tokenString, secondCommentID, postID, replyesID string) []DeleteCommentReplyesTestCase {
	return []DeleteCommentReplyesTestCase{
		{
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 200,
			TokenGiven: tokenString,
		},
		{
			//an id that does not exist CommentID
			CommentID:  strconv.Itoa(12322),
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 404,
			TokenGiven: tokenString,
		},
		{
			//an id that does not exist ReplyesID
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  strconv.Itoa(12322),
			StatusCode: 404,
			TokenGiven: tokenString,
		},
		{
			//When the user is not authenticated
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 401,
			TokenGiven: "",
		},
		{
			//When wrong token is passed
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 401,
			TokenGiven: "this is a wrong token",
		},
		{
			// When id passed is invalid CommentID
			CommentID:  "unknwon",
			PostID:     postID,
			ReplyesID:  replyesID,
			StatusCode: 400,
			TokenGiven: tokenString,
		},
		{
			// When id passed is invalid ReplyesID
			CommentID:  secondCommentID,
			PostID:     postID,
			ReplyesID:  "unknwon",
			StatusCode: 400,
			TokenGiven: tokenString,
		},
	}
}
