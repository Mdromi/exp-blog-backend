package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Mdromi/exp-blog-backend/api/mailer"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/security"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) ForgotPassword(c *gin.Context) {
	// remove any possible error, because the frontend dont reload
	errList = map[string]string{}

	// start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	user.Prepare()
	errorMessages := user.Validate("forgotpassword")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", user.Email).Take(&user).Error
	if err != nil {
		errList["No_email"] = "Sorry, we do not recognize this email"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	resetPassword := models.ResetPassword{}
	resetPassword.Prepare()

	// generate the token
	token := security.TokenHash(user.Email)
	resetPassword.Email = user.Email
	resetPassword.Token = token

	resetDetails, err := resetPassword.SvaeDatails(server.DB)

	if err != nil {
		errList = formaterror.FormatError(err.Error())
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	fmt.Println("THIS OCCURRED HERE")

	// send welcome mail to the user
	response, err := mailer.SendMail.SendResetPassword(resetDetails.Email, os.Getenv("SENDGRID_FROM"), resetDetails.Token, os.Getenv("SENDGRID_API_KEY"), os.Getenv("APP_ENV"))
	fmt.Println("response", response)
	if err != nil {
		// formattedError := formaterror.FormatError(err.Error())
		// errList = formattedError
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": response.RespBody,
	})
}

func (server *Server) ResetPassword(c *gin.Context) {
	// remove any possible error, because the frontend dont reload
	errList = map[string]string{}

	// start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	user := models.User{}
	resetPassword := models.ResetPassword{}

	err = server.DB.Debug().Model(models.ResetPassword{}).Where("token = ?", requestBody["token"]).Take(&resetPassword).Error
	if err != nil {
		errList["Invalid_token"] = "Invalid link. Try requesting again"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	if requestBody["new_password"] == "" || requestBody["retype_password"] == "" {
		errList["Empty_passwords"] = "Please ensure both field are entered"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}
	if requestBody["new_password"] != "" && requestBody["retype_password"] != "" {
		// also check if the new password
		if len(requestBody["new_password"]) < 6 || len(requestBody["retype_password"]) < 6 {
			errList["Invalid_password"] = "password should be atleast 6 characters"
			handleError(c, http.StatusUnprocessableEntity, errList)
			return
		}
	}
	if requestBody["new_password"] != requestBody["retype_password"] {
		errList["Password_unequal"] = "Passwords provided do not match"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// Note this password will be hashed before it is saved in the model
	user.Password = requestBody["new_password"]
	user.Email = resetPassword.Email

	// update the password
	user.Prepare()
	err = user.UpdatePassword(server.DB)
	if err != nil {
		fmt.Println("this is the error: ", err)
		errList["Cannot_save"] = "Cannot Save, Pls try again later"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// Delete the toke record so is not used again:
	_, err = resetPassword.DeleteDatails(server.DB)
	if err != nil {
		errList["Cannot_delete"] = "Cannot Delete record, Pls try again later"
		handleError(c, http.StatusNotFound, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Success",
	})
}
