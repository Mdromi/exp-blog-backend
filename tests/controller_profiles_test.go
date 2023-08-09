package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Mdromi/exp-blog-backend/api/controllers"
	"github.com/Mdromi/exp-blog-backend/api/models"
	utils_test_controllers "github.com/Mdromi/exp-blog-backend/api/utils/tests/controllers"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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

	samples := testdata.CreateProfileSamples(loginUserID)
	for _, v := range samples {
		// executecreateProfileTestCase(t, v, loginUserID)
		utils_test_controllers.ExecuteCreateProfileTestCase(t, v, loginUserID, &server)
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

			utils_test_controllers.AssertErrorMessages(t, responseMap, map[string]string{
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
			utils_test_controllers.AssertDeleteProfileErrorResponses(t, responseInterface)
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

	samples := testdata.UpdateProfileSamples(profileID, tokenString, loginUserID)

	for _, v := range samples {
		utils_test_controllers.ExecuteUpdateProfileTest(t, v, &server)
	}
}

func prepareTestEnvironment(t *testing.T, db *gorm.DB) {
	t.Helper()
	err := refreshUserProfileTable()
	if err != nil {
		t.Fatal(err)
	}
}
