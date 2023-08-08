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
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func (server *Server) CreateUserProfile(c *gin.Context) {
	// Clear previous error if any
	errList := map[string]string{}

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
		fmt.Println("Test 2")
		fmt.Println("Unmarshal error:", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	userModel := models.User{}
	// Check if UserID is valid and associated with an existing user
	user, err := userModel.FindUserByID(server.DB, uint32(profile.UserID))
	if err != nil {
		errList["Unauthorized"] = "Invalid UserID or user does not exist"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// TASK: Also check the user are login or not?

	if user.ProfileID != 0 {
		errList["Profile_created"] = "You already created a profile"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// Check if name, title, and bio fields are provided
	if profile.Name == "" || profile.Title == "" || profile.Bio == "" {
		errList["Missing_fields"] = "Name, title, and bio are required"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
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

	// Validate the profile fields
	errorMessages = validateProfileFields(&profile)
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

	// start processing the request
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		fmt.Println("STEP - 1")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	fmt.Println("body", body)
	requestBody := map[string]string{}
	// var requestBody = models.Profile{}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		fmt.Println("STEP - 2")
		fmt.Println("err", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// userID := requestBody["user_id"]
	// uid, err := strconv.ParseUint(userID, 10, 32)

	// Check if UserID is valid and associated with an existing user
	// _, err = FindUserByID(server.DB, uint32(uid))
	// if err != nil {
	// 	errList["Unauthorized_user"] = "Invalid UserID or user does not exist"
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"status": http.StatusUnauthorized,
	// 		"error":  errList,
	// 	})
	// 	return
	// }

	// Get user id from token for valid tokens
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		fmt.Println("STEP - 3")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// if the id is not the authentiacation user id
	if tokenID != 0 && tokenID != uint32(pid) {
		errList["Unauthorized"] = "Unauthorized"
		fmt.Println("STEP - 4")
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// Check if name, title, and bio fields are provided
	if requestBody["name"] == "" || requestBody["title"] == "" || requestBody["bio"] == "" {
		errList["Missing_fields"] = "Name, title, and bio are required"
		fmt.Println("STEP - 5")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// check for previous details
	formerProfile := models.Profile{}

	err = server.DB.Debug().Model(models.Profile{}).Where("id = ?", pid).Take(&formerProfile).Error
	if err != nil {
		errList["Profile_invalid"] = "The user is does not exist"
		fmt.Println("STEP - 6")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	newProfile := models.Profile{}

	// The password fields not entered, so update only the email
	newProfile.Name = requestBody["name"]
	newProfile.Title = requestBody["title"]
	newProfile.Bio = requestBody["bio"]

	fmt.Println("STEP - 7")
	newProfile.SocialLinks = GetSocialLinksFromBody(requestBody)

	newProfile.Prepare()
	errorMessages := newProfile.Validate("update")
	if len(errorMessages) > 0 {
		errList = errorMessages
		fmt.Println("STEP - 8")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	fmt.Println("newProfile", newProfile)
	// Validate the profile fields
	errorMessages = validateProfileFields(&newProfile)
	if len(errorMessages) > 0 {
		errList = errorMessages
		fmt.Println("STEP - 9")
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	fmt.Println("STEP - 10")
	updatedProfile, err := newProfile.UpdateAUserProfile(server.DB, uint32(pid))
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		fmt.Println("STEP - 11")
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}

	fmt.Println("updatedProfile", updatedProfile)

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

// handlare function
func validateProfileFields(profile *models.Profile) map[string]string {
	errList := map[string]string{}

	validate := validator.New()
	if err := validate.Struct(profile); err != nil {
		// Handle validation errors
		if _, ok := err.(*validator.InvalidValidationError); ok {
			// Handle error from the validation library itself (e.g., invalid struct)
			errList["Validation_error"] = "Invalid input data"
		} else {
			// Handle specific validation errors for each field
			for _, fieldErr := range err.(validator.ValidationErrors) {
				fieldName := fieldErr.Field()
				switch fieldName {
				case "Name":
					errList["Profile_name"] = "Name is required and should be between 2 and 50 characters"
				case "Title":
					errList["Profile_tile"] = "Title should be less than or equal to 100 characters"
				case "Bio":
					errList["Profile_bio"] = "Bio should be less than or equal to 500 characters"
					// Add more cases for other fields if needed
				}
			}
		}
	}

	return errList
}
