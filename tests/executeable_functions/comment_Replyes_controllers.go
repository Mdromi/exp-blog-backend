package executeablefunctions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func ExecuteCreateCommentReplyes(t *testing.T, samples []testdata.CreateCommentReplyesTestCase, server *controllers.Server) {
	for _, v := range samples {
		gin.SetMode(gin.TestMode)

		r := gin.Default()

		r.POST("/comments/:id", server.CreateCommentReplye)
		url := fmt.Sprintf("/comments/%d?commentID=%s", int(v.PostID), v.CommentID)
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(v.InputJSON))
		req.Header.Set("Authorization", v.TokenGiven)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 201 {
			responseMap := responseInterface["response"].(map[string]interface{})

			// Assuming v.CommentID is a string
			commentIDFloat, err := ConvertToFloat64(v.CommentID)
			if err != nil {
				t.Errorf("Cannot convert to uint: %v", err)
			}
			// commentIDFloat, err := strconv.ParseFloat(v.CommentID, 10)
			// if err != nil {
			// 	t.Errorf("Cannot convert to uint: %v", err)
			// }

			assert.Equal(t, responseMap["comment_id"], commentIDFloat)
			assert.Equal(t, responseMap["post_id"], float64(v.PostID))
			assert.Equal(t, responseMap["profile_id"], float64(v.ProfileID))
			assert.Equal(t, responseMap["body"], v.Body)
		}

		if v.StatusCode == 401 || v.StatusCode == 422 || v.StatusCode == 500 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}
	}
}

func ExecuteGetCommentReplyes(t *testing.T, samples []testdata.GetCommentReplyeTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()
		r.GET("/comments/:id", server.GetCommentReplyes)

		url := fmt.Sprintf("/comments/%s?commentID=%s", v.PostID, v.CommentID)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 200 {
			responseMap := responseInterface["response"].([]interface{})
			if err != nil {
				t.Errorf("Cannot convert to uint: %v", err)
			}
			assert.Equal(t, len(responseMap), v.ReplyesLength)
			assert.Equal(t, v.ProfileLength, 2)
		}
		if v.StatusCode == 400 || v.StatusCode == 404 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}
	}
}

func ExecuteUpdateCommentReplye(t *testing.T, samples []testdata.UpdateCommentReplyesTestCase, server *controllers.Server) {
	for _, v := range samples {
		gin.SetMode(gin.TestMode)

		r := gin.Default()

		r.PUT("/comments/:id", server.UpdateACommentReplyes)
		url := fmt.Sprintf("/comments/%s?commentID=%s&replyID=%s", v.PostID, v.CommentID, v.ReplyesID)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(v.UpdateJSON))
		req.Header.Set("Authorization", v.TokenGiven)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 201 {
			responseMap := responseInterface["response"].(map[string]interface{})

			// Assuming v.CommentID is a string
			commentIDFloat, err := ConvertToFloat64(v.CommentID)
			if err != nil {
				t.Errorf("Cannot convert to uint: %v", err)
			}

			assert.Equal(t, responseMap["comment_id"], commentIDFloat)
			assert.Equal(t, responseMap["post_id"], v.PostID)
			assert.Equal(t, responseMap["body"], v.Body)
		}

		if v.StatusCode == 401 || v.StatusCode == 422 || v.StatusCode == 500 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}
	}
}

func ExecuteDeleteCommentReplye(t *testing.T, samples []testdata.DeleteCommentReplyesTestCase, server *controllers.Server) {
	for _, v := range samples {

		r := gin.Default()
		r.DELETE("/comments/:id", server.DeleteComment)
		url := fmt.Sprintf("/comments/%s?commentID=%s&replyID=%s", v.PostID, v.CommentID, v.ReplyesID)
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		req.Header.Set("Authorization", v.TokenGiven)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 200 {
			responseMap := responseInterface["response"]
			assert.Equal(t, responseMap, "Comment deleted")
		}
		if v.StatusCode == 400 || v.StatusCode == 401 || v.StatusCode == 404 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}
	}
}
