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
	executeablefunctions "github.com/Mdromi/exp-blog-backend/tests/executeable_functions"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var ExecuteCreateProfileTestCase = executeablefunctions.ExecuteCreateProfileTestCase
var ExecuteDeleteProfileTest = executeablefunctions.ExecuteDeleteProfileTest
var ExecuteUpdateProfileTest = executeablefunctions.ExecuteUpdateProfileTest

// TestCreateUserProfile tests the creation of user profiles.
func TestCreateUserProfile(t *testing.T) {
	// t.Parallel()
	// Set up Gin in test mode and initialize database tables.
	gin.SetMode(gin.TestMode)
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	// Seed a user and retrieve their loginUserID.
	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}
	loginUserID := user.ID

	// Get test samples for creating user profiles and iterate over them.
	samples := testdata.CreateProfileSamples(loginUserID)
	ExecuteCreateProfileTestCase(t, samples, loginUserID, &server)
}

// TestGetUserProfile tests the retrieval of user profiles.
func TestGetUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Refresh database tables and seed user profiles.
	err := refreshAllTable()
	assert.NoError(t, err)
	_, err = seedUsersProfiles()
	assert.NoError(t, err)

	// Set up Gin and create an HTTP request for getting user profiles.
	r := gin.Default()
	r.GET("/profiles", server.GetUserProfiles)
	req, err := http.NewRequest(http.MethodGet, "/profiles", nil)
	assert.NoError(t, err)

	// Serve the HTTP request and parse the response.
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

// TestGetUserProfileByID tests the retrieval of a user profile by ID.
func TestGetUserProfileByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Refresh database tables and seed a user profile.
	err := refreshAllTable()
	assert.NoError(t, err)
	profile, err := seedOneUserProfile()
	assert.NoError(t, err)

	// Set up Gin and iterate over profile samples.
	r := gin.Default()
	r.GET("/profiles/:id", server.GetUserProfile)

	// Get test samples for get user profiles and iterate over them.
	samples := testdata.GetUserProfileSample(profile)

	for _, v := range samples {
		// Create an HTTP request for getting a user profile by ID.
		req, _ := http.NewRequest("GET", "/profiles/"+v.ID, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		// Parse the response body and perform assertions based on the test case.
		responseInterface := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseInterface)
		assert.NoError(t, err)

		assert.Equal(t, rr.Code, v.StatusCode)

		if v.StatusCode == 200 {
			responseMap, ok := responseInterface["response"].(map[string]interface{})
			assert.True(t, ok)

			userID := uint(responseMap["user_id"].(float64))
			assert.Equal(t, responseMap["name"], v.Name)
			assert.Equal(t, userID, v.UserID)
		}
		if v.StatusCode == 400 || v.StatusCode == 404 {
			responseMap, ok := responseInterface["error"].(map[string]interface{})
			assert.True(t, ok)

			executeablefunctions.AssertErrorMessages(t, responseMap, map[string]string{
				"Invalid_request": "Invalid Request",
				"No_user":         "No User Found",
			})
		}
	}
}

// TestDeleteProfile tests the deletion of user profiles.
func TestDeleteProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Refresh database tables and seed a user profile.
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	// Get test samples for get user profiles and iterate over them.
	samples := testdata.DeleteUserProfileSample()

	// Set up Gin and iterate over profile samples.
	r := gin.Default()
	r.DELETE("/profiles/:id", server.DeleteUserProfile)

	for _, v := range samples {
		// Prepare the test environment, seed a profile, and get a token.
		prepareTestEnvironment(t, server.DB)
		profile, tokenString := seedProfileAndSignIn(server.DB)

		// Update the sample with valid profile ID and token.
		if v.StatusCode != 400 {
			v.ID = strconv.Itoa(int(profile.ID))
		}
		if v.StatusCode != 401 && v.StatusCode != 400 {
			v.TokenGiven = tokenString
		}

		ExecuteDeleteProfileTest(t, v, r)
	}
}

// TestUpdateProfile tests the update of user profiles.
func TestUpdateProfile(t *testing.T) {
	// t.Parallel()

	// Set up Gin in test mode and initialize database tables.
	gin.SetMode(gin.TestMode)
	err := refreshAllTable()
	if err != nil {
		log.Fatal(err)
	}

	profile, tokenString := seedProfileAndSignIn(server.DB)
	profileID := strconv.Itoa(int(profile.ID))
	loginUserID := profile.UserID

	// Get test samples for updating user profiles and iterate over them.
	samples := testdata.UpdateProfileSamples(profileID, tokenString, loginUserID)
	ExecuteUpdateProfileTest(t, samples, &server)
}

// seedProfileAndSignIn seeds a user profile, signs in the associated user, and returns the profile and token string.
func seedProfileAndSignIn(db *gorm.DB) (models.Profile, string) {
	// Seed a user profile and find the associated user.
	profile, err := seedOneUserProfile()
	if err != nil {
		log.Fatal(err)
	}

	user, err := controllers.FindUserByID(db, uint32(profile.UserID))
	if err != nil {
		log.Fatal(err)
	}

	// Sign in the user and retrieve the authentication token.
	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"]
	tokenString := fmt.Sprintf("Bearer %v", token)

	return profile, tokenString
}

// prepareTestEnvironment prepares the test environment by refreshing the user profile table.
func prepareTestEnvironment(t *testing.T, db *gorm.DB) {
	t.Helper()
	err := refreshUserProfileTable()
	if err != nil {
		t.Fatal(err)
	}
}
