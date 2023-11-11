package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/auth"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

var handleError = formaterror.HandleError

func (server *Server) CreateUserProfile(c *gin.Context) {
	errList := make(map[string]string)

	// Parse form data
	err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		errList["Invalid_body"] = "Unable to parse form data"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}
	userID, err := strconv.ParseFloat(c.PostForm("userID"), 64)
	if err != nil {
		errList["Invalid_userID"] = "Invalid or missing userID"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// userID := c.PostForm("userID").(float64)
	fullName := c.PostForm("fullName")
	title := c.PostForm("title")
	about := c.PostForm("about")
	username := c.PostForm("username")

	// Create the profile struct with required fields
	profile := models.Profile{
		Name:     fullName,
		Title:    title,
		Bio:      about,
		UserID:   uint(userID),
		Username: username,
	}

	// Access optional fields if present
	profile.SocialLinks = &models.SocialLink{} // Initialize SocialLinks before using it

	profile.SocialLinks.Website = c.PostForm("website")
	profile.SocialLinks.Github = c.PostForm("github")
	profile.SocialLinks.Linkedin = c.PostForm("linkedin")
	profile.SocialLinks.Twitter = c.PostForm("twitter")

	// Check if UserID is valid and associated with an existing user
	userModel := models.User{}
	user, err := userModel.FindUserByID(server.DB, uint32(profile.UserID))
	if err != nil {
		errList["Not_Found_user"] = "Invalid UserID or user does not exist"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// Check if the user is logged in or not
	// TODO: Add your logic for checking if the user is logged in

	if user.ProfileID != 0 {
		errList["Profile_created"] = "You already created a profile"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// Check if required fields are provided
	if profile.Name == "" || profile.Title == "" || profile.Bio == "" {
		errList["Missing_fields"] = "Name, title, and bio are required"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// Prepare and validate the profile
	profile.Prepare()
	errorMessages := profile.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	errorMessages = ValidateProfileFields(&profile)
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// Upload profile pic if profilePic is present in the form data
	if file, err := c.FormFile("profilePic"); err == nil {
		// Call the uploadFile function with file data
		profilePicPath, err := server.uploadFile(c, uint32(profile.UserID), "profilePic", file)
		fmt.Println("profilePicPath", profilePicPath)
		if err != nil {
			errList["Cannot_Save_Profile_Pic"] = err.Error()
			handleError(c, http.StatusInternalServerError, errList)
			return
		}
		profile.ProfilePic = profilePicPath
	}

	// Upload cover pic if coverPic is present in the form data
	if file, err := c.FormFile("coverPic"); err == nil {
		// Call the uploadFile function with file data
		coverPicPath, err := server.uploadFile(c, uint32(profile.UserID), "coverPic", file)
		if err != nil {
			errList["Cannot_Save_Cover_Pic"] = err.Error()
			handleError(c, http.StatusInternalServerError, errList)
			return
		}
		profile.CoverPic = coverPicPath
	}
	// fmt.Println("profile", profile)
	// Save the profile
	// profileCreated, err := profile.SaveUserProfile(server.DB)
	// if err != nil {
	// 	formattedError := formaterror.FormatError(err.Error())
	// 	errList = formattedError
	// 	handleError(c, http.StatusBadRequest, errList)
	// 	return
	// }

	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": "",
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

func (server *Server) UpdateUserProfileImage(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	// Get image type from the request (profile_pic or cover_pic)
	imageType := c.Param("type")

	// Validate image type
	if imageType != "profile_pic" && imageType != "cover_pic" {
		errList["Invalid_request"] = "Invalid Image Type"
		handleError(c, http.StatusBadRequest, errList)
		return
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

	// Upload profile or cover pic based on the image type
	filePath, err := server.uploadFile(c, uint32(pid), "profilePic", nil)
	if err != nil {
		errList["Cannot_Save_Image"] = err.Error()
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	// Save the image path to the database
	profile := models.Profile{}
	if imageType == "profile_pic" {
		profile.ProfilePic = filePath
	} else if imageType == "cover_pic" {
		profile.CoverPic = filePath
	} else {
		errList["Invalid_request"] = "Invalid Image Type"
		handleError(c, http.StatusBadRequest, errList)
		return
	}
	profile.Prepare()
	updatedProfile, err := profile.UpdateAUserProfilePic(server.DB, uint32(pid), imageType)

	if err != nil {
		errList["Cannot_Save"] = "Cannot Save Image, Please try again later"
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
	errList := map[string]string{}

	profileID := c.Param("id")
	// check the user id is valid
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

	// TASK: Also check if the user is logged in or not?

	if user.ProfileID != uint32(pid) {
		errList["Not_Found_user"] = "Invalid UserID or user does not exist"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// if the id is not the authentication user id
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

	// Upload profile pic
	profilePicPath, err := server.uploadFile(c, uint32(pid), "profilePic", nil)
	if err != nil {
		errList["Cannot_Save_Profile_Pic"] = err.Error()
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	// Update user avatar_path if it's not the same as ProfilePic
	if user.AvatarPath != profilePicPath {
		user.AvatarPath = profilePicPath
		if err := server.DB.Save(&user).Error; err != nil {
			errList["Cannot_Update_Avatar_Path"] = "Cannot update user avatar path"
			handleError(c, http.StatusInternalServerError, errList)
			return
		}
	}

	// Set the profile pic path
	newProfile.ProfilePic = profilePicPath

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

	// Delete user profile uploads directory and its contents
	uploadsDir := "static/uploads/" + strconv.Itoa(int(pid))
	err = os.RemoveAll(uploadsDir)
	if err != nil {
		errList["Other_error"] = "Error deleting user uploads directory"
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

func (server *Server) uploadFile(c *gin.Context, userID uint32, fieldname string, file *multipart.FileHeader) (string, error) {
	// Base directory without the user-specific folder
	baseDir := filepath.Join("static", "uploads", fieldname)

	// Ensure the base directory exists
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return "", errors.New("could not create upload directory")
	}

	// Create user-specific directory within the base directory
	userDir := filepath.Join(baseDir, strconv.Itoa(int(userID)))
	if err := os.MkdirAll(userDir, os.ModePerm); err != nil {
		return "", errors.New("could not create user-specific upload directory")
	}

	// Save file to a new file in the user-specific directory
	filePath := filepath.Join(userDir, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.New("could not save file on server")
	}

	return filePath, nil
}
