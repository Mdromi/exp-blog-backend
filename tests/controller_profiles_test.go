package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type createProfileTestCase struct {
	inputJSON  string
	statusCode int
	userID     uint
	name       string
	title      string
	profilePic string
	socialLink models.SocialLink
}

func TestCreateUserProfile(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	loginUserID := user.ID

	type SocialLink struct {
		Facebook  string
		Twitter   string
		Instagram string
	}

	samples := []createProfileTestCase{
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com",
					"twitter": "www.twitter.com",
					"instagram": "www.instagram.com"
				}
			}`, loginUserID),
			statusCode: 201,
			userID:     loginUserID,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, 0),
			statusCode: 404,
			userID:     0,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, 342049902),
			statusCode: 404,
			userID:     0,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			statusCode: 400,
			userID:     loginUserID,
			name:       "Pet",
			title:      "",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			statusCode: 400,
			userID:     loginUserID,
			name:       "",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, user.ID),
			statusCode: 400,
			userID:     loginUserID,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "", 
				"title": "", 
				"bio": "", 
				"profile_pic": "/images/profile.jpg", "social_links": {}
			}`, loginUserID),
			statusCode: 400,
			userID:     loginUserID,
			name:       "",
			title:      "",
			profilePic: "/images/profile.jpg",
		},
		{
			inputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "", "title": "This is the title", "bio": "This is the Bio", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			statusCode: http.StatusBadRequest,
		},
		{
			inputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "Pet", "title": "", "bio": "This is the Bio", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			statusCode: http.StatusBadRequest,
		},
		{
			inputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "Pet", "title": "This is the title", "bio": "", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			statusCode: http.StatusBadRequest,
		},
		{
			inputJSON:  fmt.Sprintf(`{ "user_id": %d, "name": "", "title": "", "bio": "", "profile_pic": "/images/profile.jpg", "social_links": {} }`, loginUserID),
			statusCode: http.StatusBadRequest,
		},
	}

	for _, v := range samples {
		executecreateProfileTestCase(t, v, loginUserID)
	}
}

func executecreateProfileTestCase(t *testing.T, v createProfileTestCase, loginUserID uint) {
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

func TestGetUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	assert.NoError(t, err)

	_, err = seedUsersProfiles()
	assert.NoError(t, err)

	r := gin.Default()
	r.GET("/profiles", server.GetUserProfiles)

	req, err := http.NewRequest(http.MethodGet, "/profiles", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	profilesMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &profilesMap)
	assert.NoError(t, err)

	// Check the response code and the number of user profiles
	theUserProfiles, ok := profilesMap["response"].([]interface{})
	assert.True(t, ok)
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theUserProfiles), 2)
}

func TestGetUserProfileByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	assert.NoError(t, err)

	profile, err := seedOneUserProfile()
	assert.NoError(t, err)

	profileSample := []struct {
		id         string
		statusCode int
		name       string
		userID     uint
	}{
		{
			id:         strconv.Itoa(int(profile.ID)),
			statusCode: 200,
			name:       profile.Name,
			userID:     profile.UserID,
		},
		{
			id:         "unknwon",
			statusCode: 400,
		},
		{
			id:         strconv.Itoa(12322), // an id that does not exist
			statusCode: 404,
		},
	}

	r := gin.Default()
	r.GET("/profiles/:id", server.GetUserProfile)

	for _, v := range profileSample {
		req, _ := http.NewRequest("GET", "/profiles/"+v.id, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		assert.NoError(t, err)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap, ok := responseInterface["response"].(map[string]interface{})
			assert.True(t, ok)

			userID := uint(responseMap["user_id"].(float64))
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, userID, v.userID)
		}
		if v.statusCode == 400 || v.statusCode == 404 {
			responseMap, ok := responseInterface["error"].(map[string]interface{})
			assert.True(t, ok)

			assertErrorMessages(t, responseMap, map[string]string{
				"Invalid_request": "Invalid Request",
				"No_user":         "No User Found",
			})
		}
	}
}

func TestDeleteProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	profileSample := []struct {
		id         string
		tokenGiven string
		statusCode int
	}{
		{
			id:         "",
			tokenGiven: "",
			statusCode: 200,
		},
		{
			id:         "",
			tokenGiven: "",
			statusCode: 401,
		},
		{
			id:         "",
			tokenGiven: "This is an incorrect token",
			statusCode: 401,
		},
		{
			id:         "unknown",
			tokenGiven: "",
			statusCode: 400,
		},
	}

	r := gin.Default()
	r.DELETE("/profiles/:id", server.DeleteUserProfile)

	for _, v := range profileSample {
		prepareTestEnvironment(t, server.DB)

		profile, tokenString := seedProfileAndSignIn(server.DB)

		if v.statusCode != 400 {
			v.id = strconv.Itoa(int(profile.ID))
		}
		if v.statusCode != 401 && v.statusCode != 400 {
			v.tokenGiven = tokenString
		}

		req, _ := http.NewRequest(http.MethodDelete, "/profiles/"+v.id, nil)
		req.Header.Set("Authorization", v.tokenGiven)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err := json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, responseInterface["response"], "User deleted")
		}
		if v.statusCode == 400 || v.statusCode == 401 {
			assertDeleteProfileErrorResponses(t, responseInterface)
		}
	}
}

func seedProfileAndSignIn(db *gorm.DB) (models.Profile, string) {
	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatal(err)
	}

	user, err := controllers.FindUserByID(db, uint32(profile.UserID))
	if err != nil {
		log.Fatal(err)
	}

	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"]
	tokenString := fmt.Sprintf("Bearer %v", token)

	return profile, tokenString
}

func assertDeleteProfileErrorResponses(t *testing.T, responseInterface map[string]interface{}) {
	responseMap := responseInterface["error"].(map[string]interface{})

	errorMessages := map[string]string{
		"Invalid_request": "Invalid Request",
		"Unauthorized":    "Unauthorized",
	}

	assertErrorMessages(t, responseMap, errorMessages)
}

type updateProfileTestCase struct {
	ID         string
	UpdateJSON string
	UserID     uint
	StatusCode int
	Name       string
	Title      string
	ProfilePic string
	TokenGiven string
}

func TestUpdateProfile(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatal(err)
	}

	// profileID := profile.ID
	profileID := strconv.Itoa(int(profile.ID))

	loginUserID := profile.UserID

	// Check if UserID is valid and associated with an existing user
	user, err := controllers.FindUserByID(server.DB, uint32(loginUserID))
	if err != nil {
		log.Fatal(err)
	}

	// Login the user and get the authentication token
	tokenInterface, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	token := tokenInterface["token"] // get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []updateProfileTestCase{

		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d, 
				"name": "Pet 1", 
				"title": "This is the title - 1", 
				"bio": "This is the Bio - 1", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com/mdromi",
					"twitter": "www.twitter.com/mdromi",
					"instagram": "www.instagram.com/mdromi"
				}
			}`, loginUserID),
			StatusCode: 200,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, 0),
			StatusCode: 401,
			UserID:     0,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: "342049902",
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 404,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "",
				"title": "This is the title",
				"bio": "This is the Bio",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "Pet",
				"title": "This is the title",
				"bio": "",
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "Pet",
			Title:      "This is the title",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "",
				"title": "",
				"bio": "",
				"profile_pic": "/images/profile.jpg", "social_links": {}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "",
			Title:      "",
			ProfilePic: "/images/profile.jpg",
			TokenGiven: tokenString,
		},
		{
			ID: profileID,
			UpdateJSON: fmt.Sprintf(`{
				"user_id": %d,
				"name": "",
				"title": "",
				"bio": "",
				"profile_pic": "/images/profile.jpg", "social_links": {}
			}`, loginUserID),
			StatusCode: 422,
			UserID:     loginUserID,
			Name:       "",
			Title:      "",
			TokenGiven: tokenString,
		},
	}

	for _, v := range samples {
		executeUpdateProfileTest(t, v)
	}
}

func executeUpdateProfileTest(t *testing.T, v updateProfileTestCase) {
	r := gin.Default()
	r.PUT("/profiles/:id", server.UpdateAUserProfile)

	req, err := http.NewRequest(http.MethodPut, "/profiles/"+v.ID, bytes.NewBufferString(v.UpdateJSON))
	if err != nil {
		t.Fatalf("Error creating request: %v\n", err)
	}

	req.Header.Set("Authorization", v.TokenGiven)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	var responseInterface map[string]interface{}
	err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
	if err != nil {
		t.Errorf("Cannot convert to JSON: %v", err)
	}

	assert.Equal(t, rr.Code, v.StatusCode)

	if v.StatusCode == http.StatusOK {
		responseMap := responseInterface["response"].(map[string]interface{})
		assert.Equal(t, responseMap["name"], "Pet 1")
		assert.Equal(t, responseMap["title"], "This is the title - 1")
		assert.Equal(t, responseMap["profile_pic"], "image/pic")
		assert.Equal(t, responseMap["user_id"], float64(v.UserID))
	} else {
		assertErrorResponse(t, responseInterface["error"].(map[string]interface{}), v.StatusCode)
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

func prepareTestEnvironment(t *testing.T, db *gorm.DB) {
	t.Helper()
	err := refreshUserProfileTable()
	if err != nil {
		t.Fatal(err)
	}
}
