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
	"github.com/Mdromi/exp-blog-backend/api/utils/fileformat"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

var handleError = formaterror.HandleError

func (server *Server) CreateUserProfile(c *gin.Context) {
	// Clear previous error if any
	errList := map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	profile := models.Profile{}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	userModel := models.User{}
	// Check if UserID is valid and associated with an existing user
	user, err := userModel.FindUserByID(server.DB, uint32(profile.UserID))
	if err != nil {
		errList["Not_Found_user"] = "Invalid UserID or user does not exist"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// TASK: Also check the user are login or not?

	if user.ProfileID != 0 {
		errList["Profile_created"] = "You already created a profile"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// Check if name, title, and bio fields are provided
	if profile.Name == "" || profile.Title == "" || profile.Bio == "" {
		errList["Missing_fields"] = "Name, title, and bio are required"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	profile.Prepare()
	errorMessages := profile.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// Validate the profile fields
	errorMessages = ValidateProfileFields(&profile)
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	profileCreated, err := profile.SaveUserProfile(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		handleError(c, http.StatusBadRequest, errList)
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
		handleError(c, http.StatusInternalServerError, errList)
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
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	profile := models.Profile{}

	profileGotten, err := profile.FindUserProfileByID(server.DB, uint32(pid))
	if err != nil {
		errList["No_profile"] = "No Profile Found"
		handleError(c, http.StatusNotFound, errList)
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
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// Get user id from the token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// if the id is not the authenticated user id
	if tokenID != 0 && tokenID != uint32(pid) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	f, err := file.Open()
	if err != nil {
		errList["Invalid_file"] = "Invalid File"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}
	defer f.Close()

	size := file.Size
	// The image should not be more than 500KB
	if size > int64(512000) {
		errList["To_large"] = "Sorry, Please upload an Image of 500KB or less"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	buffer := make([]byte, size)
	f.Read(buffer)
	fileBytes := bytes.NewReader(buffer)
	fileType := http.DetectContentType(buffer)
	// if the image is valid
	if !strings.HasPrefix(fileType, "image") {
		errList["Not_Image"] = "Please Upload a valid image"
		handleError(c, http.StatusUnprocessableEntity, errList)
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
		handleError(c, http.StatusInternalServerError, errList)
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
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// Get user id from token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// Check if profile is valid and associated with an existing user
	profile, err := FindUserProfileByID(server.DB, uint32(pid))
	if err != nil {
		errList["Not_Found_profile"] = "Not Found the profile"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// Check if UserID is valid and associated with an existing user
	user, err := FindUserByID(server.DB, uint32(profile.UserID))
	if err != nil {
		errList["Not_Found_user"] = "Invalid UserID or user does not exist"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// TASK: Also check the user are login or not?

	if user.ProfileID != uint32(pid) {
		errList["Not_Found_user"] = "Invalid UserID or user does not exist"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// if the id is not the authentiacation user id
	if tokenID != 0 && tokenID != uint32(pid) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	newProfile := models.Profile{}

	err = json.Unmarshal(body, &newProfile)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// Check if name, title, and bio fields are provided
	if newProfile.Name == "" || newProfile.Title == "" || newProfile.Bio == "" {
		errList["Missing_fields"] = "Name, title, and bio are required"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	newProfile.Prepare()
	errorMessages := newProfile.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		if errorMessages["user_id"] != "" {
			handleError(c, http.StatusUnauthorized, errList)
			return
		}
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}
	// Validate the profile fields
	errorMessages = ValidateProfileFields(&newProfile)
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	updatedProfile, err := newProfile.UpdateAUserProfile(server.DB, uint32(pid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		handleError(c, http.StatusInternalServerError, errList)
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
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// Check if profile is valid and associated with an existing user
	profile, err := FindUserProfileByID(server.DB, uint32(pid))
	if err != nil {
		errList["Not_Found_profile"] = "Not Found the profile"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// get user id from the token for valid tokens
	tokenID, err = auth.ExtractTokenID(c.Request)
	fmt.Println("tokenID", tokenID)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// If the id is not the authenticated user id
	if tokenID != 0 && tokenID != uint32(profile.UserID) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// Also delete the posts, likes and the comments that this user created if any:

	comment := models.Comment{}
	likeDislike := models.LikeDislike{}
	post := models.Post{}

	_, err = post.DeleteUserPosts(server.DB, uint32(pid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusInternalServerError, errList)
		return
	}
	_, err = comment.DeleteUserComments(server.DB, uint32(pid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusInternalServerError, errList)
		return
	}
	_, err = likeDislike.DeleteUserLikes(server.DB, uint32(pid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	deletedProfile := models.Profile{}
	_, err = deletedProfile.DeleteAUserProfile(server.DB, uint32(pid))
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "User deleted",
	})
}
