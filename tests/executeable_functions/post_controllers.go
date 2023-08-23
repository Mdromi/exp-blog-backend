package executeablefunctions

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func ExecuteCreatePostTest(t *testing.T, samples []testdata.CreatePostTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()

		r.POST("/posts", server.CreatePost)
		req, err := http.NewRequest(http.MethodPost, "/posts", bytes.NewBufferString(v.InputJSON))
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
			assert.Equal(t, responseMap["title"], v.Title)
			assert.Equal(t, responseMap["content"], v.Content)
			assert.Equal(t, responseMap["tags"], v.Tags)
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

func ExecuteGetPostByID(t *testing.T, samples []testdata.GetPostByIdTestCase, server *controllers.Server) {
	for _, v := range samples {
		req, _ := http.NewRequest("GET", "/posts/"+v.ID, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/posts/:id", server.GetPost)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err := json.Unmarshal(rr.Body.Bytes(), &responseInterface)

		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.StatusCode)
		if v.StatusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["title"], v.Title)
			assert.Equal(t, responseMap["content"], v.Content)
			assert.Equal(t, responseMap["author_id"], float64(v.Author_id))
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

func ExecuteUpdatePost(t *testing.T, samples []testdata.UpdatePostTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()

		r.PUT("/posts/:id", server.UpdatePost)
		req, err := http.NewRequest(http.MethodPut, "/posts/"+v.ID, bytes.NewBufferString(v.UpdateJSON))
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

		if v.StatusCode == 200 {
			//casting the interface to map:
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["title"], v.Title)
			assert.Equal(t, responseMap["content"], v.Content)
		}
		if v.StatusCode == 400 || v.StatusCode == 401 || v.StatusCode == 422 || v.StatusCode == 500 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}

	}
}

func ExecuteDeletePost(t *testing.T, samples []testdata.DeletePostTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()
		r.DELETE("/posts/:id", server.DeletePost)
		req, _ := http.NewRequest(http.MethodDelete, "/posts/"+v.ID, nil)
		req.Header.Set("Authorization", v.TokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})

		err := json.Unmarshal(rr.Body.Bytes(), &responseInterface)

		if err != nil {
			t.Errorf("Cannot convert to json here: %v", err)
		}
		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 200 {
			assert.Equal(t, responseInterface["response"], "Post deleted")
		}

		if v.StatusCode == 400 || v.StatusCode == 401 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}
	}
}
