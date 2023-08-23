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

func (server *Server) CreateCommentReplye(c *gin.Context) {

	// Extrect CommentID
	// POST /posts/123?commentID=123
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

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	replye := models.Replyes{}
	err = json.Unmarshal(body, &replye)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// erter the profile, comment and the postid. the reply body is automatically passed
	replye.ProfileID = uint64(profileID)
	replye.PostID = uint32(pid)
	replye.CommentID = cid

	replye.Preapre()
	errorMessages := replye.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	commentReplyCreated, err := replye.SaveCommentReplyes(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		handleError(c, http.StatusNotFound, errList)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": commentReplyCreated,
	})
}

func (server *Server) GetCommentReplyes(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")
	_, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

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

	// check if the comment exist
	origComment := models.Comment{}

	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&origComment).Error

	if err != nil {
		errList["No_comment"] = "No Comment Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	replye := models.Replyes{}
	replyes, err := replye.GetCommentReplyes(server.DB, cid)
	if err != nil {
		errList["No_comment_replyes"] = "No Comment Replyes Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": replyes,
	})
}

func (server *Server) UpdateACommentReplyes(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	pid, profileID, user, post := server.CommonCommentAndReplyesCode(c)
	if pid == 0 || profileID == 0 || user == nil || post == nil {
		return
	}

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

	replyID := c.Query("replyID")
	if replyID == "" {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}
	rcid, err := strconv.ParseUint(replyID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// check if the comment exist
	origComment := models.Comment{}

	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&origComment).Error

	if err != nil {
		errList["No_comment"] = "No Comment Found"
		handleError(c, http.StatusNotFound, errList)
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	// check if the comment replyes exist
	origCommentReplyes := models.Replyes{}
	err = server.DB.Debug().Model(models.Replyes{}).Where("id = ? AND comment_id = ? AND profile_id = ?", rcid, cid, profileID).Take(&origCommentReplyes).Error
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	if profileID != uint32(origCommentReplyes.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// start processing requested data
	replye := models.Replyes{}
	err = json.Unmarshal(body, &replye)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	replye.Preapre()
	errorMessages := replye.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		handleError(c, http.StatusUnprocessableEntity, errList)
		return
	}

	replye.ID = origCommentReplyes.ID //this is important to tell the model the post id to update, the other update field are set above
	replye.ProfileID = origCommentReplyes.ProfileID
	replye.PostID = origCommentReplyes.PostID
	replye.CommentID = origCommentReplyes.CommentID
	replye.PostID = origCommentReplyes.PostID

	commentReplyUpdated, err := replye.UpdateACommentReplyes(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		handleError(c, http.StatusInternalServerError, errList)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": commentReplyUpdated,
	})
}

func (server *Server) DeleteCommentReplye(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	pid, profileID, user, post := server.CommonCommentAndReplyesCode(c)
	if pid == 0 || profileID == 0 || user == nil || post == nil {
		return
	}
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

	replyID := c.Query("replyID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}
	rcid, err := strconv.ParseUint(replyID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	// check if the comment replyes exist
	origCommentReplyes := models.Replyes{}
	err = server.DB.Debug().Model(models.Replyes{}).Where("id = ? AND comment_id = ? AND profile_id = ?", rcid, cid, profileID).Take(&origCommentReplyes).Error
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		handleError(c, http.StatusBadRequest, errList)
		return
	}

	if profileID != uint32(origCommentReplyes.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// Is the authenticated user, the owner of this replye?
	if profileID != uint32(origCommentReplyes.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		handleError(c, http.StatusUnauthorized, errList)
		return
	}

	// If all the conditions are met, delete the post
	_, err = origCommentReplyes.DeleteAReplyes(server.DB)
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
