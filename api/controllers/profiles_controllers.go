package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Mdromi/exp-blog-backend/api/auth"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/security"
	"github.com/Mdromi/exp-blog-backend/api/utils/fileformat"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) CreateUserProfile(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	profile := models.Profile{}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	profile.Prepare()
	errorMessages := profile.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	profileCreated, err := profile.SaveUserProfile(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": profileCreated,
	})
}

func (server *Server) GetUserProfiles(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	profile := models.Profile{}
	profiles, err := profile.FindAllUsersProfile(server.DB)
	if err != nil {
		errList["No_profile"] = "No Profile Found"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": profiles,
	})
}

func (server *Server) GetUserProfile(c *gin.Context) {
	profileId := c.Param("id")

	pid, err := strconv.ParseUint(profileId, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	profile := models.Profile{}

	profileGotten, err := profile.FindUserProfileByID(server.DB, uint32(pid))
	if err != nil {
		errList["No_profile"] = "No Profile Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": profileGotten,
	})
}

func (server *Server) UpdateUserProfilePic(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error gotting env, %v", err)
	}
	profileId := c.Param("id")
	// check if the user id is valid
	pid, err := strconv.ParseUint(profileId, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Get user id from the token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// if the id is not the authenticated user id
	if tokenID != 0 && tokenID != uint32(pid) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	defer f.Close()

	size := file.Size
	// The image should not be more than 500KB
	if size > int64(512000) {
		errList["To_large"] = "Sorry, Please upload an Image of 500KB or less"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	buffer := make([]byte, size)
	f.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	// if the image is valid
	if !strings.HasPrefix(fileType, "image") {
		errList["Not_Image"] = "Please Upload a valid image"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	filePath := fileformat.UniqueFormat(file.Filename)
	path := "/profile-photos/" + filePath
	params := &s3.PutObjectInput{
		Bucket:        aws.String("chodapi"),
		Key:           aws.String(path),
		Body:          fileBytes,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(fileType),
		ACL:           aws.String("public-read"),
	}
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("DO_SPACES_KEY"), os.Getenv("DO_SPACES_SECRET"), os.Getenv("DO_SPACES_TOKEN")),
		Endpoint: aws.String(os.Getenv("DO_SPACES_ENDPOINT")),
		Region:   aws.String(os.Getenv("DO_SPACES_REGION")),
	}
	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	_, err = s3Client.PutObject(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//IF YOU PREFER TO USE AMAZON S3
	// s, err := session.NewSession(&aws.Config{Too_large
	// 	Region: aws.String("us-east-1"),
	// 	Credentials: credentials.NewStaticCredentials(
	// 		os.Getenv("AWS_KEY"),
	// 		os.Getenv("AWS_SECRET"),
	// 		os.Getenv("AWS_TOKEN"),
	// 		),
	// })
	// if err != nil {
	// 	fmt.Printf("Could not upload file first error: %s\n", err)
	// }
	// fileName, err := SaveProfileImage(s, file)
	// if err != nil {
	// 	fmt.Printf("Could not upload file %s\n", err)
	// } else {
	// 	fmt.Printf("Image uploaded: %s\n", fileName)
	// }

	// save The iamge path to the database
	profile := models.Profile{}
	profile.ProfilePic = filePath
	profile.Prepare()
	updatedProfile, err := profile.UpdateAUserProfilePic(server.DB, uint32(pid))

	if err != nil {
		errList["Cannot_Save"] = "Cannot Save Image, Pls try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedProfile,
	})
}

// TASK: NEED TO MODIFIED
func (server *Server) UpdateAUserProfile(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	profileID := c.Param("id")
	// check the user id is  valid
	pid, err := strconv.ParseUint(profileID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Get user id from token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// if the id is not the authentiacation user id
	if tokenID != 0 && tokenID != uint32(pid) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	requestBody := map[string]string{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// check for previous details
	formerProfile := models.Profile{}

	err = server.DB.Debug().Model(models.Profile{}).Where("id = ?", pid).Take(&formerProfile).Error
	if err != nil {
		errList["Profile_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// NOTE: USING PROFILE CONTROLLER UNDER USER
	formerUser := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", formerProfile.UserID).Take(&formerUser).Error
	if err != nil {
		errList["User_invalid"] = "The user is does not exist"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newProfile := models.Profile{}
	// NEXT : UPGRADE THIS FUNCTIONALITY
	newUser := models.User{}
	// when current password has content.
	if requestBody["current_password"] == "" && requestBody["new_password"] != "" {
		errList["Empty_current"] = "Please Provide current password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	if requestBody["current_password"] != "" && requestBody["new_password"] == "" {
		errList["Empty_current"] = "Please Provide current password"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	if requestBody["current_password"] != "" && requestBody["new_password"] != "" {
		// also check if the new password
		if len(requestBody["new_password"]) < 6 {
			errList["Invalid_password"] = "Password should be atleast 6 characters"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		// if they do, check that the former password is correct
		err = security.VerifyPassword(formerUser.Password, requestBody["current_password"])

		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			errList["Password_mismatch"] = "The password not correct"
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"status": http.StatusUnprocessableEntity,
				"error":  errList,
			})
			return
		}

		// update both the password and the email
		newUser.Email = requestBody["email"]
		newUser.Password = requestBody["new_password"]
	}

	// The password fields not entered, so update only the email
	newProfile.Name = requestBody["name"]
	newProfile.Title = requestBody["title"]
	newProfile.Bio = requestBody["bio"]

	newUser.Prepare()
	errorMessages := newUser.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newProfile.Prepare()
	errorMessages = newProfile.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	_, err = newUser.UpdateAUser(server.DB, uint32(formerProfile.UserID))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	updatedProfile, err := newProfile.UpdateAUserProfile(server.DB, uint32(pid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": updatedProfile,
	})
}

func (server *Server) DeleteUserProfile(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}
	var tokenID uint32
	profileID := c.Param("id")

	// check if the user id is valid
	pid, err := strconv.ParseUint(profileID, 10, 32)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// get user id from the token for valid tokens
	tokenID, err = auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If the id is not the authenticated user id
	if tokenID != 0 && tokenID != uint32(pid) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	profile := models.Profile{}
	_, err = profile.DeleteAUserProfile(server.DB, uint32(pid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Also delete the posts, likes and the comments that this user created if any:

	// comment := models.Comment{}
	// like := models.Like{}
	// post := models.Post{}

	// _, err = post.DeleteUserPosts(server.DB, uint32(pid))
	// if err != nil {
	// 	errList["Other_error"] = "Please try again later"
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": http.StatusInternalServerError,
	// 		"error":  err,
	// 	})
	// 	return
	// }
	// _, err = comment.DeleteUserComments(server.DB, uint32(pid))
	// if err != nil {
	// 	errList["Other_error"] = "Please try again later"
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": http.StatusInternalServerError,
	// 		"error":  err,
	// 	})
	// 	return
	// }
	// _, err = like.DeleteUserLikes(server.DB, uint32(pid))
	// if err != nil {
	// 	errList["Other_error"] = "Please try again later"
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": http.StatusInternalServerError,
	// 		"error":  err,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User deleted",
	})
}
