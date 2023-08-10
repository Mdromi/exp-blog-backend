package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	executeablefunctions "github.com/Mdromi/exp-blog-backend/tests/executeable_functions"
	"github.com/Mdromi/exp-blog-backend/tests/testdata"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	// Get test samples for creating user and iterate over them.
	samples := testdata.CreateUserSamples()
	executeablefunctions.ExecuteCreateUserTestCase(t, samples, &server)
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	_, err = seedUsers()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/users", server.GetUsers)

	req, err := http.NewRequest(http.MethodGet, "/users", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	usersMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(rr.Body.String()), &usersMap)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}

	// This is so that we can get the length of the users
	theUsers := usersMap["response"].([]interface{})
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(theUsers), 2)
}

func TestGetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	// Get test samples for getUserByID  and iterate over them.
	samples := testdata.GetUserByIDSamples(user)
	executeablefunctions.ExecuteGetUserByIdTestCase(t, samples, &server)
}

func TestUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var AuthEmail, AuthPassword, AuthUsername, AuthEmail2 string
	var AuthID uint32

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	users, err := seedUsers() // we need atleast two users to properly check the update
	if err != nil {
		log.Fatalf("Error seeding user: %v\n", err)
	}

	// Get only the first user
	for _, user := range users {
		if user.ID == 2 {
			AuthEmail2 = user.Email
			continue
		}
		AuthID = uint32(user.ID)
		AuthEmail = user.Email
		AuthUsername = user.Username
		AuthPassword = "password" //Note the password in the database is already hashed, we want unhashed
	}

	// Login the user and get the authentication token
	tokenInterface, err := server.SignIn(AuthEmail, AuthPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}

	token := tokenInterface["token"] // get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get test samples for update user  and iterate over them.
	samples := testdata.UpdateUserSamples(tokenString, AuthUsername, AuthEmail2, AuthID)
	executeablefunctions.ExecuteUpdateUserTest(t, samples, &server)
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		log.Fatal(err)
	}

	// Note: The value of the user password before it was hashed is "password". so:
	password := "password"
	tokenInterface, err := server.SignIn(user.Email, password)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	token := tokenInterface["token"] // get only the token
	tokenString := fmt.Sprintf("Bearer %v", token)

	userID := strconv.Itoa(int(user.ID))

	// Get test samples for delete user  and iterate over them.
	samples := testdata.DeleteUserSample(tokenString, userID)
	executeablefunctions.ExecuteDeleteUserTest(t, samples, &server)
}
