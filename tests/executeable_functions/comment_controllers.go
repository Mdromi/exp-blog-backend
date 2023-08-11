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

func ExecuteCreateComments(t *testing.T, samples []testdata.CreateCommentTestCase, server *controllers.Server) {
	for _, v := range samples {
		gin.SetMode(gin.TestMode)

		r := gin.Default()

		r.POST("/comments/:id", server.CreateComment)
		req, err := http.NewRequest(http.MethodPost, "/comments/"+v.PostIDString, bytes.NewBufferString(v.InputJSON))
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

func ExecuteGetComments(t *testing.T, samples []testdata.GetCommentTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()
		r.GET("/comments/:id", server.GetComments)
		req, err := http.NewRequest(http.MethodGet, "/comments/"+v.PostID, nil)
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
			assert.Equal(t, len(responseMap), v.CommentsLength)
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

func ExecuteUpdateComments(t *testing.T, samples []testdata.UpdateCommentsTestCase, server *controllers.Server, postID float64, secondUserID float64) {
	for _, v := range samples {
		r := gin.Default()
		r.PUT("/comments/:id", server.UpdateComment)
		url := fmt.Sprintf("/comments/%d?commentID=%s", int(postID), v.CommentID)
		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBufferString(v.UpdateJSON))
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
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["post_id"], postID)
			assert.Equal(t, responseMap["profile_id"], secondUserID)
			assert.Equal(t, responseMap["body"], v.Body)
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

func ExecuteDeleteComments(t *testing.T, samples []testdata.DeleteCommentsTestCase, server *controllers.Server, postID float64) {
	for _, v := range samples {

		r := gin.Default()
		r.DELETE("/comments/:id", server.DeleteComment)
		// req, err := http.NewRequest(http.MethodDelete, "/comments/"+v.CommentID, nil)
		url := fmt.Sprintf("/comments/%d?commentID=%s", int(postID), v.CommentID)
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
