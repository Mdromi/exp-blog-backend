package executeablefunctions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/Mdromi/exp-blog-backend/tests/testdata"
)

// ExecuteCreateProfileTestCase executes a test case for creating a user profile.
func ExecuteCreateProfileTestCase(t *testing.T, samples []testdata.CreateProfileTestCase, loginUserID uint, server *controllers.Server) {
	for _, v := range samples {
		// Set up the Gin router and route for creating a user profile.
		r := gin.Default()
		r.POST("/profiles", server.CreateUserProfile)

		// Create an HTTP request for the profile creation endpoint.
		req, err := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(v.InputJSON))
		if err != nil {
			t.Fatalf("this is the error: %v\n", err)
		}

		// Serve the HTTP request and record the response.
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// Parse the response body into a map.
		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}

		// Check if the expected status code matches the actual response.
		if v.StatusCode == http.StatusCreated {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["name"], v.Name)
			assert.Equal(t, responseMap["user_id"], float64(v.UserID))
		} else {
			errorResponse, ok := responseInterface["error"].(map[string]interface{})
			if !ok {
				t.Errorf("Received unexpected response format: %v", responseInterface)
			} else {
				AssertErrorResponse(t, errorResponse, v.StatusCode)
			}
		}

		// Clean up: reset the user's profile_id to 0 in the database.
		err = server.DB.Debug().Model(&models.User{}).Where("id = ?", loginUserID).Update("profile_id", 0).Error
		if err != nil {
			fmt.Println("err", err)
			log.Fatal(err)
		}
	}
}

// ExecuteUpdateProfileTest executes a test case for updating a user profile.
func ExecuteUpdateProfileTest(t *testing.T, samples []testdata.UpdateProfileTestCase, server *controllers.Server) {
	for _, v := range samples {
		// Set up the Gin router and route for updating a user profile.
		r := gin.Default()
		r.PUT("/profiles/:id", server.UpdateAUserProfile)

		// Create an HTTP request for the profile update endpoint.
		req, err := http.NewRequest(http.MethodPut, "/profiles/"+v.ID, bytes.NewBufferString(v.UpdateJSON))
		if err != nil {
			t.Fatalf("Error creating request: %v\n", err)
		}

		// Set the Authorization header and serve the HTTP request.
		req.Header.Set("Authorization", v.TokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// Parse the response body into a map.
		var responseInterface map[string]interface{}
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to JSON: %v", err)
		}

		// Check if the response status code matches the expected status code.
		assert.Equal(t, rr.Code, v.StatusCode)

		// Check the response details if the status code is http.StatusOK.
		if v.StatusCode == http.StatusOK {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["name"], "Pet 1")
			assert.Equal(t, responseMap["title"], "This is the title - 1")
			assert.Equal(t, responseMap["profile_pic"], "image/pic")
			assert.Equal(t, responseMap["user_id"], float64(v.UserID))
		} else {
			AssertErrorResponse(t, responseInterface["error"].(map[string]interface{}), v.StatusCode)
		}
	}
}

func ExecuteDeleteProfileTest(t *testing.T, v testdata.DeleteUserProfileSampleCase, r *gin.Engine) {
	// Create an HTTP request for deleting a user profile.
	req, _ := http.NewRequest(http.MethodDelete, "/profiles/"+v.ID, nil)
	req.Header.Set("Authorization", v.TokenGiven)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Parse the response body and perform assertions.
	responseInterface := make(map[string]interface{})
	err := json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}
	assert.Equal(t, rr.Code, v.StatusCode)

	if v.StatusCode == 200 {
		assert.Equal(t, responseInterface["response"], "User deleted")
	}
	if v.StatusCode == 400 || v.StatusCode == 401 {
		AssertDeleteProfileErrorResponses(t, responseInterface)
	}
}
