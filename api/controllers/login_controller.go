package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Mdromi/exp-blog-backend/api/auth"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/security"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errorMessage := map[string]string{"first error": "Unable to get request"}
		handleError(c, http.StatusUnprocessableEntity, errorMessage)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		// Convert the untyped string constant into a map[string]string
		errorMessage := map[string]string{"error": "Cannot unmarshal body"}
		handleError(c, http.StatusUnprocessableEntity, errorMessage)
		return
	}
	errorMessages := user.Validate("login")
	if len(errorMessages) > 0 {
		handleError(c, http.StatusUnprocessableEntity, errorMessages)
		return
	}

	userData, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		handleError(c, http.StatusUnprocessableEntity, formattedError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": userData,
	})
}

func (server *Server) SignIn(email, password string) (map[string]interface{}, error) {
	var err error

	userData := make(map[string]interface{})

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		fmt.Println("this is the error getting the user:", err)
		return nil, err
	}

	err = security.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("this is the error hashing the password: ", err)
		return nil, err
	}
	token, err := auth.CreateToken(uint32(user.ID))
	if err != nil {
		fmt.Println("this is the error creating the token: ", err)
		return nil, err
	}

	userData["token"] = token
	userData["id"] = user.ID
	userData["email"] = user.Email
	userData["avatar_path"] = user.AvatarPath
	userData["username"] = user.Username

	return userData, nil
}
