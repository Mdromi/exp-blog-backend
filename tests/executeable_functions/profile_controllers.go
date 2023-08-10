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

// AssertErrorResponse asserts error responses based on the provided status code.
func AssertErrorResponse(t *testing.T, responseMap map[string]interface{}, statusCode int) {
	// Define error messages for various status codes.
	errorMessages := map[int]map[string]string{
		http.StatusBadRequest: {
			"Unmarshal_error": "Cannot unmarshal body",
			"Profile_name":    "Name is required and should be between 2 and 50 characters",
			"Profile_title":   "Title should be less than or equal to 100 characters",
			"Profile_bio":     "Bio should be less than or equal to 500 characters",
			"Missing_fields":  "Name, title, and bio are required",
			"Required_name":   "Name is required",
			"Name_length":     "Name must be between 2 and 50 characters",
			"Required_title":  "Title is required",
			"Title_length":    "Title cannot exceed 100 characters",
			"Required_bio":    "Bio is required",
			"Bio_length":      "Bio cannot exceed 500 characters",
		},
		http.StatusUnauthorized: {
			"Unauthorized":      "Unauthorized",
			"Unauthorized_user": "Invalid UserID or user does not exist",
		},
		http.StatusNotFound: {
			"Not_Found_profile": "Not Found the profile",
			"Not_Found_user":    "Invalid UserID or user does not exist",
		},
		http.StatusInternalServerError: {
			"Internal_error": "Internal server error occurred",
			// You can add more error messages specific to http.StatusInternalServerError here...
		},
		http.StatusUnprocessableEntity: {
			"Unprocessable_entity": "Request could not be processed",
			"Unmarshal_error":      "Cannot unmarshal body",
			"Profile_created":      "You already created a profile",
			"Invalid_body":         "Unable to get request",
			// Define error messages for http.StatusUnprocessableEntity here...
		},
		// ... add more status code error messages as needed ...
	}

	// Check if error messages are defined for the given status code.
	if errorMsgs, ok := errorMessages[statusCode]; ok {
		AssertErrorMessages(t, responseMap, errorMsgs)
	} else {
		t.Errorf("No error messages defined for status code: %d", statusCode)
	}
}

// AssertErrorMessages asserts specific error messages in the response map.
func AssertErrorMessages(t *testing.T, responseMap map[string]interface{}, errorMessages map[string]string) {
	for key, expected := range errorMessages {
		if responseMap[key] != nil {
			fmt.Println("statusCode, errorMsgs", responseMap[key])
			assert.Equal(t, responseMap[key], expected)
		}
	}
}

// AssertDeleteProfileErrorResponses asserts error responses for profile deletion.
func AssertDeleteProfileErrorResponses(t *testing.T, responseInterface map[string]interface{}) {
	responseMap := responseInterface["error"].(map[string]interface{})

	errorMessages := map[string]string{
		"Invalid_request": "Invalid Request",
		"Unauthorized":    "Unauthorized",
	}

	AssertErrorMessages(t, responseMap, errorMessages)
}
