package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) CreateComment(c *gin.Context) {
	pid, profileID, user, post := server.CommonCommentAndReplyesCode(c)
	if pid == 0 || profileID == 0 || user == nil || post == nil {
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// erter the profile and the postid. the comment body is automatically passed
	comment.ProfileID = profileID
	comment.PostID = pid

	comment.Preapre()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	commentCreated, err := comment.SaveComment(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		handleError(c, http.StatusNotFound, errList)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": commentCreated,
	})
}

func (server *Server) GetComments(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")

	// Is a valdi post id given to us?
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// check if the post exist:
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	comment := models.Comment{}

	comments, err := comment.GetComments(server.DB, pid)
	if err != nil {
		errList["No_comments"] = "No comments found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": comments,
	})
}

func (server *Server) UpdateComment(c *gin.Context) {

	// PUT /comment/123?commentID=102
	commentID := c.Query("commentID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}
	cid, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	pid, profileID, user, post := server.CommonCommentAndReplyesCode(c)
	if pid == 0 || profileID == 0 || user == nil || post == nil {
		return
	}

	origComment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&origComment).Error
	if err != nil {
		errList["No_comment"] = "No Comment Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	if profileID != origComment.ProfileID {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// start processing requested data
	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	comment.Preapre()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	comment.ID = origComment.ID //this is important to tell the model the post id to update, the other update field are set above
	comment.ProfileID = origComment.ProfileID
	comment.PostID = origComment.PostID

	commentUpdated, err := comment.UpdateAComment(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		handleError(c, http.StatusInternalServerError, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": commentUpdated,
	})
}

func (server *Server) DeleteComment(c *gin.Context) {
	// DELETE /comment/123?commentID=102
	commentID := c.Query("commentID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}
	cid, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	pid, profileID, user, post := server.CommonCommentAndReplyesCode(c)
	if pid == 0 || profileID == 0 || user == nil || post == nil {
		return
	}

	origComment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&origComment).Error
	if err != nil {
		errList["No_comment"] = "No Comment Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	if profileID != origComment.ProfileID {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// If all the conditions are met, delete the post
	_, err = origComment.DeleteAComment(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		handleError(c, http.StatusNotFound, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Comment deleted",
	})
}
