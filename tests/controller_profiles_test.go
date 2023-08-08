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
)

func TestCreateUserProfile(t *testing.T) {
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

	// SocialLink represents social links for a user's profile
	type SocialLink struct {
		Facebook  string
		Twitter   string
		Instagram string
	}

	samples := []struct {
		inputJSON  string
		statusCode int
		userID     uint
		name       string
		title      string
		profilePic string
		socialLink SocialLink
	}{
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
			statusCode: 401,
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
			statusCode: 401,
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
		r := gin.Default()
		r.POST("/profiles", server.CreateUserProfile)
		req, err := http.NewRequest(http.MethodPost, "/profiles", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		err = server.DB.Debug().Model(&models.User{}).Where("id = ?", loginUserID).Update("profile_id", 0).Error
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, responseMap["user_id"], uint64(v.userID))
		}
		if v.statusCode == 422 || v.statusCode == 500 || v.statusCode == 400 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_body"] != nil {
				assert.Equal(t, responseMap["Invalid_body"], "Unable to get request")
			}
			if responseMap["Unmarshal_error"] != nil {
				assert.Equal(t, responseMap["Unmarshal_error"], "Cannot unmarshal body")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Invalid UserID or user does not exist")
			}
			if responseMap["Profile_created"] != nil {
				assert.Equal(t, responseMap["Profile_created"], "You already created a profile")
			}
			if responseMap["Profile_name"] != nil {
				assert.Equal(t, responseMap["Profile_name"], "Name is required and should be between 2 and 50 characters")
			}
			if responseMap["Profile_tile"] != nil {
				assert.Equal(t, responseMap["Profile_tile"], "Title should be less than or equal to 100 characters")
			}
			if responseMap["Profile_bio"] != nil {
				assert.Equal(t, responseMap["Profile_bio"], "Bio should be less than or equal to 500 characters")
			}
			if responseMap["Missing_fields"] != nil {
				assert.Equal(t, responseMap["Missing_fields"], "Name, title, and bio are required")
			}
		}
	}
}

func TestGetUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsersProfiles()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/profiles", server.GetUserProfiles)

	req, err := http.NewRequest(http.MethodGet, "/profiles", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	profilesMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(rr.Body.String()), &profilesMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	// This is so that we can get the length of the users
	theUserProfiles := profilesMap["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theUserProfiles), 2)
}

func TestGetUserProfileByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatal(err)
	}

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
			id:         strconv.Itoa(12322), //an id that does not exist
			statusCode: 404,
		},
	}

	for _, v := range profileSample {
		req, _ := http.NewRequest("GET", "/profiles/"+v.id, nil)
		rr := httptest.NewRecorder()

		r := gin.Default()
		r.GET("/profiles/:id", server.GetUserProfile)
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v\n", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			userID := uint(responseMap["user_id"].(float64))
			assert.Equal(t, responseMap["name"], v.name)
			assert.Equal(t, userID, v.userID)
		}
		if v.statusCode == 400 || v.statusCode == 404 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_request"] != nil {
				assert.Equal(t, responseMap["Invalid_request"], "Invalid Request")
			}

			if responseMap["No_user"] != nil {
				assert.Equal(t, responseMap["No_user"], "No User Found")
			}
		}
	}
}

func TestUpdateProfile(t *testing.T) {
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

	loginUserID := strconv.Itoa(int(profile.UserID))

	// Check if UserID is valid and associated with an existing user
	user, err := controllers.FindUserByID(server.DB, uint32(profile.UserID))
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

	samples := []struct {
		id         string
		updateJSON string
		userID     string
		statusCode int
		name       string
		title      string
		bio        string
		profilePic string
		tokenGiven string
	}{

		{
			id: profileID,
			updateJSON: fmt.Sprintf(`{
				"user_id": "%s", 
				"name": "Pet 1", 
				"title": "This is the title - 1", 
				"bio": "This is the Bio - 1", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com/mdromi",
					"twitter": "www.twitter.com/mdromi",
					"instagram": "www.instagram.com/mdromi"
				}
			}`, loginUserID),
			statusCode: 201,
			userID:     loginUserID,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
			tokenGiven: tokenString,
		},
		{
			id: "",
			updateJSON: fmt.Sprintf(`{
				"user_id":"%s", 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			statusCode: 401,
			userID:     "0",
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
			tokenGiven: tokenString,
		},
		{
			id: "342049902",
			updateJSON: fmt.Sprintf(`{
				"user_id": "%s", 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "This is the Bio", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			statusCode: 401,
			userID:     loginUserID,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
			tokenGiven: tokenString,
		},
		{
			id: profileID,
			updateJSON: fmt.Sprintf(`{
				"user_id": "%s", 
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
			tokenGiven: tokenString,
		},
		{
			id: profileID,
			updateJSON: fmt.Sprintf(`{
				"user_id": "%s", 
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
			tokenGiven: tokenString,
		},
		{
			id: profileID,
			updateJSON: fmt.Sprintf(`{
				"user_id": "%s", 
				"name": "Pet", 
				"title": "This is the title", 
				"bio": "", 
				"profile_pic": "/images/profile.jpg", "social_links": {
					"facebook": "www.facebook.com", "twitter": "www.twitter.com", "instagram": "www.instagram.com"
				}
			}`, loginUserID),
			statusCode: 400,
			userID:     loginUserID,
			name:       "Pet",
			title:      "This is the title",
			profilePic: "/images/profile.jpg",
			tokenGiven: tokenString,
		},
		{
			id: profileID,
			updateJSON: fmt.Sprintf(`{
				"user_id": "%s", 
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
			tokenGiven: tokenString,
		},
		{
			id: profileID,
			updateJSON: fmt.Sprintf(`{ 
				"user_id": "%s", 
				"name": "", 
				"title": "", 
				"bio": "", 
				"profile_pic": "/images/profile.jpg", "social_links": {} 
			}`, loginUserID),
			statusCode: http.StatusBadRequest,
			tokenGiven: tokenString,
		},
	}

	for i, v := range samples {
		r := gin.Default()
		r.PUT("/profiles/:id", server.UpdateAUserProfile)

		req, err := http.NewRequest(http.MethodPut, "/profiles/"+v.id, bytes.NewBufferString(v.updateJSON))
		req.Header.Set("Authorization", v.tokenGiven)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		responseInterface := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseInterface)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		// err = server.DB.Debug().Model(&models.User{}).Where("id = ?", loginUserID).Update("profile_id", 0).Error
		// if err != nil {
		// 	log.Fatal(err)
		// }
		assert.Equal(t, rr.Code, v.statusCode)
		fmt.Println("Itaratot", i)

		if v.statusCode == 200 {
			responseMap := responseInterface["response"].(map[string]interface{})
			assert.Equal(t, responseMap["name"], "Pet 1")
			assert.Equal(t, responseMap["title"], "This is the title - 1")
			assert.Equal(t, responseMap["bio"], "This is the Bio - 1")
			assert.Equal(t, responseMap["user_id"], loginUserID)
		}

		if v.statusCode == 422 || v.statusCode == 500 || v.statusCode == 400 {
			responseMap := responseInterface["error"].(map[string]interface{})

			if responseMap["Invalid_body"] != nil {
				assert.Equal(t, responseMap["Invalid_body"], "Unable to get request")
			}
			if responseMap["Unmarshal_error"] != nil {
				assert.Equal(t, responseMap["Unmarshal_error"], "Cannot unmarshal body")
			}
			if responseMap["Unauthorized"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Unauthorized")
			}
			if responseMap["Unauthorized_user"] != nil {
				assert.Equal(t, responseMap["Unauthorized"], "Invalid UserID or user does not exist")
			}
			if responseMap["Profile_created"] != nil {
				assert.Equal(t, responseMap["Profile_created"], "You already created a profile")
			}
			if responseMap["Profile_name"] != nil {
				assert.Equal(t, responseMap["Profile_name"], "Name is required and should be between 2 and 50 characters")
			}
			if responseMap["Profile_tile"] != nil {
				assert.Equal(t, responseMap["Profile_tile"], "Title should be less than or equal to 100 characters")
			}
			if responseMap["Profile_bio"] != nil {
				assert.Equal(t, responseMap["Profile_bio"], "Bio should be less than or equal to 500 characters")
			}
			if responseMap["Missing_fields"] != nil {
				assert.Equal(t, responseMap["Missing_fields"], "Name, title, and bio are required")
			}
		}
	}
}
