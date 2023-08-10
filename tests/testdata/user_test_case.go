package testdata

import (
	"fmt"
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/models"
)

type CreateUserTestCase struct {
	InputJSON  string
	StatusCode int
	Username   string
	Email      string
}

type GetUserByIDTestCase struct {
	ID         string
	StatusCode int
	Username   string
	Email      string
}

type UpdateUserTestCase struct {
	ID          string
	UpdateJSON  string
	StatusCode  int
	Username    string
	UpdateEmail string
	TokenGiven  string
}

type DeleteUserTestCase struct {
	ID         string
	TokenGiven string
	StatusCode int
}

func CreateUserSamples() []CreateUserTestCase {
	return []CreateUserTestCase{
		{
			InputJSON:  `{"username":"Pet", "email": "pet@example.com", "password": "password"}`,
			StatusCode: 201,
			Username:   "Pet",
			Email:      "pet@example.com",
		},
		{
			InputJSON:  `{"username":"Frank", "email": "pet@example.com", "password": "password"}`,
			StatusCode: 500,
		},
		{
			InputJSON:  `{"username":"Pet", "email": "grand@example.com", "password": "password"}`,
			StatusCode: 500,
		},
		{
			InputJSON:  `{"username":"Kan", "email": "kanexample.com", "password": "password"}`,
			StatusCode: 422,
		},
		{
			InputJSON:  `{"username": "", "email": "kan@example.com", "password": "password"}`,
			StatusCode: 422,
		},
		{
			InputJSON:  `{"username": "Kan", "email": "", "password": "password"}`,
			StatusCode: 422,
		},
		{
			InputJSON:  `{"username": "Kan", "email": "kan@example.com", "password": ""}`,
			StatusCode: 422,
		},
	}
}

func GetUserByIDSamples(user models.User) []GetUserByIDTestCase {
	return []GetUserByIDTestCase{
		{
			ID:         strconv.Itoa(int(user.ID)),
			StatusCode: 200,
			Username:   user.Username,
			Email:      user.Email,
		},
		{
			ID:         "unknwon",
			StatusCode: 400,
		},
		{
			ID:         strconv.Itoa(12322), //an id that does not exist
			StatusCode: 404,
		},
	}
}

func UpdateUserSamples(tokenString, AuthUsername, AuthEmail2 string, AuthID uint32) []UpdateUserTestCase {
	return []UpdateUserTestCase{
		{
			// In this particular test case, we changed the user's password to "newpassword". Very important to note
			// Convert int32 to int first before converting to string
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"email": "grand@example.com", "current_password": "password", "new_password": "newpassword"}`,
			StatusCode:  200,
			Username:    AuthUsername, //the username does not change, even if a new name is provided, it will be ignored
			UpdateEmail: "grand@example.com",
			TokenGiven:  tokenString,
		},
		{
			// An attempt to change the username, will not work, the old name is still retained.
			// Remember, the "current_password" is now "newpassword", changed in test 1
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"username": "new_name", "email": "grand@example.com", "current_password": "newpassword", "new_password": "newpassword"}`,
			StatusCode:  200,
			Username:    AuthUsername, //irrespective of the username inputed above, the old one is still used
			UpdateEmail: "grand@example.com",
			TokenGiven:  tokenString,
		},
		{
			// The user can update only his email address
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"email": "fred@example.com"}`,
			StatusCode:  200,
			Username:    AuthUsername,
			UpdateEmail: "fred@example.com",
			TokenGiven:  tokenString,
		},
		{
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"email": "alex@example.com", "current_password": "", "new_password": ""}`,
			StatusCode:  200,
			Username:    AuthUsername,
			UpdateEmail: "alex@example.com",
			TokenGiven:  tokenString,
		},
		{
			// When password the "current_password" is given and does not match with the one in the database
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"email": "alex@example.com", "current_password": "wrongpassword", "new_password": "password"}`,
			StatusCode:  422,
			UpdateEmail: "alex@example.com",
			TokenGiven:  tokenString,
		},
		{
			// When password the "current_password" is correct but the "new_password" field is not given
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"email": "alex@example.com", "current_password": "newpassword", "new_password": ""}`,
			StatusCode:  422,
			UpdateEmail: "alex@example.com",
			TokenGiven:  tokenString,
		},
		{
			// When password the "current_password" is correct but the "new_password" field is not up to 6 characters
			ID:          strconv.Itoa(int(AuthID)),
			UpdateJSON:  `{"email": "alex@example.com", "current_password": "newpassword", "new_password": "pass"}`,
			StatusCode:  422,
			UpdateEmail: "alex@example.com",
			TokenGiven:  tokenString,
		},
		{
			// When no token was passed (when the user is not authenticated)
			ID:         strconv.Itoa(int(AuthID)),
			UpdateJSON: `{"email": "man@example.com", "current_password": "newpassword", "new_password": "password"}`,
			StatusCode: 401,
			TokenGiven: "",
		},
		{
			// When incorrect token was passed
			ID:         strconv.Itoa(int(AuthID)),
			UpdateJSON: `{"email": "man@example.com", "current_password": "newpassword", "new_password": "password"}`,
			StatusCode: 401,
			TokenGiven: "This is incorrect token",
		},
		{
			// Remember "kenny@example.com" belongs to user 2, so, user 1 cannot use some else email that is in our database
			ID:         strconv.Itoa(int(AuthID)),
			UpdateJSON: fmt.Sprintf(`{"email": "%s", "current_password": "newpassword", "new_password": "password"}`, AuthEmail2),
			StatusCode: 500,
			TokenGiven: tokenString,
		},
		{
			// When the email provided is invalid
			ID:         strconv.Itoa(int(AuthID)),
			UpdateJSON: `{"email": "notexample.com", "current_password": "newpassword", "new_password": "password"}`,
			StatusCode: 422,
			TokenGiven: tokenString,
		},
		{
			// If the email field is empty
			ID:         strconv.Itoa(int(AuthID)),
			UpdateJSON: `{"email": "", "current_password": "newpassword", "new_password": "password"}`,
			StatusCode: 422,
			TokenGiven: tokenString,
		},
		{
			// when invalid is provided
			ID:         "unknwon",
			TokenGiven: tokenString,
			StatusCode: 400,
		},
	}
}

func DeleteUserSample(tokenString string, userID string) []DeleteUserTestCase {
	return []DeleteUserTestCase{
		{
			// Convert int32 to int first before converting to string
			ID:         userID,
			TokenGiven: tokenString,
			StatusCode: 200,
		},
		{
			// When no token is given
			ID:         userID,
			TokenGiven: "",
			StatusCode: 401,
		},
		{
			// When incorrect token is given
			ID:         userID,
			TokenGiven: "This is an incorrect token",
			StatusCode: 401,
		},
		{
			// When bad request data is given:
			ID:         "unknwon",
			TokenGiven: tokenString,
			StatusCode: 400,
		},
	}
}
