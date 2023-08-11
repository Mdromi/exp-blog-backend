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

func ExecuteCreateUserTestCase(t *testing.T, samples []testdata.CreateUserTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()
		r.POST("/users", server.CreateUser)
		req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(v.InputJSON))
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
			// casting the interface to map:
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["username"], v.Username)
			assert.Equal(t, responseMap["email"], v.Email)
		}

		if v.StatusCode == 422 || v.StatusCode == 500 {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}
	}
}

func ExecuteGetUserByIdTestCase(t *testing.T, samples []testdata.GetUserByIDTestCase, server *controllers.Server) {
	for _, v := range samples {
		req, _ := http.NewRequest("GET", "/users/"+v.ID, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/users/:id", server.GetUser)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err := json.Unmarshal(rr.Body.Bytes(), &responseInterface)

		if err != nil {
			t.Errorf("Cannot convert to json: %v\n", err)
		}

		assert.Equal(t, rr.Code, v.StatusCode)
		if v.StatusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["username"], v.Username)
			assert.Equal(t, responseMap["email"], v.Email)
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

func ExecuteUpdateUserTest(t *testing.T, samples []testdata.UpdateUserTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()

		r.PUT("/users/:id", server.UpdateUser)
		req, err := http.NewRequest(http.MethodPut, "/users/"+v.ID, bytes.NewBufferString(v.UpdateJSON))
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
			// casting the interface to map:
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["email"], v.UpdateEmail)
			assert.Equal(t, responseMap["username"], v.Username)
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

func ExecuteDeleteUserTest(t *testing.T, samples []testdata.DeleteUserTestCase, server *controllers.Server) {
	for _, v := range samples {
		r := gin.Default()
		r.DELETE("/users/:id", server.DeleteUser)
		req, _ := http.NewRequest(http.MethodDelete, "/users/"+v.ID, nil)
		req.Header.Set("Authorization", v.TokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err := json.Unmarshal(rr.Body.Bytes(), &responseInterface)

		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 200 {
			assert.Equal(t, responseInterface["response"], "User deleted")
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
