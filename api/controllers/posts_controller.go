package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/auth"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/Mdromi/exp-blog-backend/api/utils/postformator"
	"github.com/gin-gonic/gin"
)

func (server *Server) CreatePost(c *gin.Context) {
	// cleat previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	post := models.Post{}

	err = json.Unmarshal(body, &post)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	pid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}
	// check if the user exist:
	profile := models.Profile{}
	err = server.DB.Debug().Model(models.Profile{}).Where("id = ?", pid).Take(&profile).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	post.AuthorID = uint(pid) // the authenticated user is the one creating the post

	post.Prepare()
	errorMessages := post.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	if len(post.Tags) == 0 {
		errList["Invalid_tags"] = "Invalid Tagas"
		handleError(c, http.StatusBadRequest, errList)
		return
	}
	// result := ConvertTags(post.Tags)

	// postPermalinks := postformator.CreatePostPermalinks(post.Title)

	post.PostPermalinks = postformator.CreatePostPermalinks(post.Title)
	post.ReadTime = postformator.CalculateReadingTime(post.Content)

	postCreated, err := post.SavePost(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		handleError(c, http.StatusInternalServerError, errList)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": postCreated,
	})

}

func (server *Server) GetPosts(c *gin.Context) {
	post := models.Post{}

	posts, err := post.FindAllPosts(server.DB)
	if err != nil {
		errList["No_post"] = "No Post Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": posts,
	})
}

func (server *Server) GetPost(c *gin.Context) {
	postID := c.Param("id")
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	post := models.Post{}
	postReceived, err := post.FindPostById(server.DB, pid)
	if err != nil {
		errList["No_post"] = "No Post Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": postReceived,
	})
}

func (server *Server) UpdatePost(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")
	// check if the post id is valid
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	//Check if the auth token is valid and  get the user id from it
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	//Check if the post exist
	origPost := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&origPost).Error

	if err != nil {
		errList["No_post"] = "No Post Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}
	if profileID != uint32(origPost.AuthorID) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// Read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// start processing the request data
	post := models.Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	post.ID = origPost.ID // this is important to tell the model the post id to update, the other field are set above
	post.AuthorID = origPost.AuthorID

	post.Prepare()
	errorMessages := post.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	post.PostPermalinks = postformator.CreatePostPermalinks(post.Title)
	post.ReadTime = postformator.CalculateReadingTime(post.Content)

	postUpdated, err := post.UpdateAPost(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		handleError(c, http.StatusInternalServerError, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": postUpdated,
	})
}

func (server *Server) DeletePost(c *gin.Context) {
	postID := c.Param("id")
	// is a valid post id given to us?
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	fmt.Println("this is delete post sir")

	// is this user authenticated?
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// check if the post exist
	post := models.Post{}
	err = server.DB.Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// Is the authenticated user, the owner of this post?
	if profileID != uint32(post.AuthorID) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// if all the conditions are metn delete the post
	_, err = post.DeleteAPost(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	commnnt := models.Comment{}
	likeDislike := models.LikeDislike{}

	// also delete the likes and the comments that thi post have:
	_, err = commnnt.DeletePostComments(server.DB, pid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	_, err = likeDislike.DeletePostLikes(server.DB, pid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusInternalServerError, errList)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Post deleted",
	})
}

func (server *Server) GetUserProfilePosts(c *gin.Context) {
	profileID := c.Param("id")
	// IS a valid user id given to us ?
	pid, err := strconv.ParseUint(profileID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	post := models.Post{}
	posts, err := post.FindUserPosts(server.DB, uint32(pid))

	if err != nil {
		errList["No_post"] = "No Post Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": posts,
	})
}
