package executeablefunctions

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertErrorResponse asserts error responses based on the provided status code.
func AssertErrorResponse(t *testing.T, responseMap map[string]interface{}, statusCode int) {
	// Define error messages for various status codes.
	errorMessages := map[int]map[string]string{
		http.StatusBadRequest: {
			"Unmarshal_error":  "Cannot unmarshal body",
			"Profile_name":     "Name is required and should be between 2 and 50 characters",
			"Profile_title":    "Title should be less than or equal to 100 characters",
			"Profile_bio":      "Bio should be less than or equal to 500 characters",
			"Missing_fields":   "Name, title, and bio are required",
			"Required_name":    "Name is required",
			"Name_length":      "Name must be between 2 and 50 characters",
			"Required_title":   "Title is required",
			"Title_length":     "Title cannot exceed 100 characters",
			"Required_bio":     "Bio is required",
			"Bio_length":       "Bio cannot exceed 500 characters",
			"Invalid_tags":     "Invalid_tags",
			"Taken_title":      "Title Already Taken",
			"Required_content": "Required Content",
			"Required_author":  "Required Author",
			"Invalid_request":  "Invalid Request",
		},
		http.StatusUnauthorized: {
			"Unauthorized":      "Unauthorized",
			"Unauthorized_user": "Invalid UserID or user does not exist",
		},
		http.StatusNotFound: {
			"Not_Found_profile":  "Not Found the profile",
			"Not_Found_user":     "Invalid UserID or user does not exist",
			"No_post":            "No Post Found",
			"No_comment":         "No Comment Found",
			"No_comment_replyes": "No Comment Replyes Found",
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
			"Required_body":        "Required Comment",
			"No_email":             "Sorry, we do not recognize this email",
			"Invalid_token":        "Invalid link. Try requesting again",
			"Empty_passwords":      "Please ensure both field are entered",
			"Password_unequal":     "Passwords provided do not match",
			"Cannot_save":          "Cannot Save, Pls try again later",
			"Cannot_delete":        "Cannot Delete record, Pls try again later",
			"Required_email":       "Required Email",
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

func ConvertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("unsupported data type")
	}
}
