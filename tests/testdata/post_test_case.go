package testdata

import (
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/models"
)

type CreatePostTestCase struct {
	InputJSON  string
	StatusCode int
	Title      string
	Content    string
	TokenGiven string
	Tags       []interface{}
}

type GetPostByIdTestCase struct {
	ID         string
	StatusCode int
	Title      string
	Content    string
	Author_id  uint32
}

type UpdatePostTestCase struct {
	ID         string
	UpdateJSON string
	StatusCode int
	Title      string
	Content    string
	TokenGiven string
}

type DeletePostTestCase struct {
	ID           string
	TokenGiven   string
	StatusCode   int
	ErrorMessage string
}

func CreatePostsSamples(tokenString string) []CreatePostTestCase {
	return []CreatePostTestCase{
		{
			InputJSON:  `{"title":"The title", "content": "the content", "tags": ["tag1", "tag2", "tag3"], "thumbnails": "img/thumbnails.png"}`,
			StatusCode: 201,
			TokenGiven: tokenString,
			Title:      "The title",
			Content:    "the content",
			Tags:       []interface{}{"tag1", "tag2", "tag3"},
		},
		{
			// When the post title already exist
			InputJSON:  `{"title":"The title", "content": "the content", "tags": ["tag1", "tag2", "tag3"]}`,
			StatusCode: 500,
			TokenGiven: tokenString,
		},
		{
			// When no token is passed
			InputJSON:  `{"title":"When no token is passed", "content": "the content", "tags": ["tag1", "tag2", "tag3"]}`,
			StatusCode: 401,
			TokenGiven: "",
		},
		{
			// When incorrect token is passed
			InputJSON:  `{"title":"When incorrect token is passed", "content": "the content"}`,
			StatusCode: 401,
			TokenGiven: "This is an incorrect token",
		},
		{
			InputJSON:  `{"title": "", "content": "The content"}`,
			StatusCode: 400,
			TokenGiven: tokenString,
		},
		{
			InputJSON:  `{"title": "This is a title", "content": ""}`,
			StatusCode: 400,
			TokenGiven: tokenString,
		},
	}
}

func GetPostByIDSamples(post models.Post) []GetPostByIdTestCase {
	return []GetPostByIdTestCase{
		{
			ID:         strconv.Itoa(int(post.ID)),
			StatusCode: 200,
			Title:      post.Title,
			Content:    post.Content,
			Author_id:  uint32(post.AuthorID),
		},
		{
			ID:         "unknwon",
			StatusCode: 400,
		},
		{
			ID:         strconv.Itoa(12322), //an id that does not exist
			StatusCode: 404,
		},
	}
}

func UpdatePostTestSamples(tokenString string, AuthPostID uint) []UpdatePostTestCase {
	return []UpdatePostTestCase{
		{
			// Convert int64 to int first before converting to string
			ID:         strconv.Itoa(int(AuthPostID)),
			UpdateJSON: `{"title":"The updated post", "content": "This is the updated content", "tags": ["tag1", "tag2", "tag3"], "thumbnails": "img/thumbnails.png"}`,
			StatusCode: 200,
			Title:      "The updated post",
			Content:    "This is the updated content",
			TokenGiven: tokenString,
		},
		{
			// When no token is provided
			ID:         strconv.Itoa(int(AuthPostID)),
			UpdateJSON: `{"title":"This is still another title", "content": "This is the updated content"}`,
			TokenGiven: "",
			StatusCode: 401,
		},
		{
			// When incorrect token is provided
			ID:         strconv.Itoa(int(AuthPostID)),
			UpdateJSON: `{"title":"This is still another title", "content": "This is the updated content"}`,
			TokenGiven: "this is an incorrect token",
			StatusCode: 401,
		},
		{
			//Note: "Title 2" belongs to post 2, and title must be unique
			ID:         strconv.Itoa(int(AuthPostID)),
			UpdateJSON: `{"title":"Title 2", "content": "This is the updated content"}`,
			StatusCode: 500,
			TokenGiven: tokenString,
		},
		{
			// When title is not given
			ID:         strconv.Itoa(int(AuthPostID)),
			UpdateJSON: `{"title":"", "content": "This is the updated content"}`,
			StatusCode: 422,
			TokenGiven: tokenString,
		},
		{
			// When content is not given
			ID:         strconv.Itoa(int(AuthPostID)),
			UpdateJSON: `{"title":"Awesome title", "content": ""}`,
			StatusCode: 422,
			TokenGiven: tokenString,
		},
		{
			// When invalid post id is given
			ID:         "unknwon",
			StatusCode: 400,
		},
	}
}

func DeletePostTestSamples(tokenString string, AuthPostID uint) []DeletePostTestCase {
	return []DeletePostTestCase{
		{
			// Convert int64 to int first before converting to string
			ID:         strconv.Itoa(int(AuthPostID)),
			TokenGiven: tokenString,
			StatusCode: 200,
		},
		{
			// When empty token is passed
			ID:         strconv.Itoa(int(AuthPostID)),
			TokenGiven: "",
			StatusCode: 401,
		},
		{
			// When incorrect token is passed
			ID:         strconv.Itoa(int(AuthPostID)),
			TokenGiven: "This is an incorrect token",
			StatusCode: 401,
		},
		{
			ID:         "unknwon",
			TokenGiven: tokenString,
			StatusCode: 400,
		},
		{
			ID:           strconv.Itoa(int(1)),
			StatusCode:   401,
			ErrorMessage: "Unauthorized",
		},
	}
}
