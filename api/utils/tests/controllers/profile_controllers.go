package utils_test_controllers

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
	"github.com/Mdromi/exp-blog-backend/api/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var createProfileTestCase = tests.CreateProfileTestCase
var server = controllers.Server{}

func ExecutecreateProfileTestCase(t *testing.T, v createProfileTestCase, loginUserID uint) {
	r := gin.Default()
	r.POST("/profiles", server.CreateUserProfile)
	req, err := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(v.inputJSON))
	if err != nil {
		t.Fatalf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	responseInterface := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
	}

	if v.statusCode == http.StatusCreated {
		responseMap := responseInterface["response"].(map[string]interface{})
		assert.Equal(t, responseMap["name"], v.name)
		assert.Equal(t, responseMap["user_id"], float64(v.userID))
	} else {
		errorResponse, ok := responseInterface["error"].(map[string]interface{})
		if !ok {
			t.Errorf("Received unexpected response format: %v", responseInterface)
		} else {
			assertErrorResponse(t, errorResponse, v.statusCode)
		}
	}

	err = server.DB.Debug().Model(&models.User{}).Where("id = ?", loginUserID).Update("profile_id", 0).Error
	if err != nil {
		fmt.Println("err", err)
		log.Fatal(err)
	}
}

func assertErrorResponse(t *testing.T, responseMap map[string]interface{}, statusCode int) {
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
			"Unauthorized":      "Invalid UserID or user does not exist",
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

	if errorMsgs, ok := errorMessages[statusCode]; ok {
		assertErrorMessages(t, responseMap, errorMsgs)
	} else {
		t.Errorf("No error messages defined for status code: %d", statusCode)
	}
}

func assertErrorMessages(t *testing.T, responseMap map[string]interface{}, errorMessages map[string]string) {
	for key, expected := range errorMessages {
		if responseMap[key] != nil {
			fmt.Println("statusCode, errorMsgs", responseMap[key])
			assert.Equal(t, responseMap[key], expected)
		}
	}
}
