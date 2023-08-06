package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/auth"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) CreateCommentReplye(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// check the token
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the profile exists;
	profile := models.Profile{}
	err = server.DB.Debug().Model(models.Profile{}).Where("id = ?", profileID).Take(&profile).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// Extrect CommentID
	// GET /posts/123?commentID=123
	commentID := c.Query("commentID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	cid, err := strconv.ParseUint(commentID, 10, 64)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	replye := models.Replyes{}
	err = json.Unmarshal(body, &replye)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	commentReplyCreated, err := replye.SaveCommentReplyes(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	commentID := c.Query("commentID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	cid, err := strconv.ParseUint(commentID, 10, 64)

	replye := models.Replyes{}
	replyes, err := replye.GetCommentReplyes(server.DB, cid)
	if err != nil {
		errList["No_comments"] = "No comments found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
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

	postID := c.Param("id")
	_, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	commentID := c.Query("commentID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	cid, err := strconv.ParseUint(commentID, 10, 64)

	replyID := c.Query("replyID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	rcid, err := strconv.ParseUint(replyID, 10, 64)

	// check if the auth token is valid and get the user id from it
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the comment exist
	origComment := models.Comment{}

	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&origComment).Error

	if err != nil {
		errList["No_comment"] = "No Comment Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	// read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// check if the comment replyes exist
	origCommentReplyes := models.Replyes{}
	err = server.DB.Debug().Model(models.Replyes{}).Where("id = ? AND comment_id = ? AND profile_id = ?", rcid, cid, profileID).Take(&origCommentReplyes).Error

	if profileID != uint32(origCommentReplyes.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// start processing requested data
	replye := models.Replyes{}
	err = json.Unmarshal(body, &replye)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	replye.Preapre()
	errorMessages := replye.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	replye.ID = origCommentReplyes.ID //this is important to tell the model the post id to update, the other update field are set above
	replye.ProfileID = origCommentReplyes.ProfileID
	replye.PostID = origCommentReplyes.PostID

	commentReplyUpdated, err := replye.UpdateACommentReplyes(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": commentReplyUpdated,
	})
}

func (server *Server) DeleteCommentReplye(c *gin.Context) {
	commentID := c.Query("commentID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	cid, err := strconv.ParseUint(commentID, 10, 64)

	// is this user authenticated?
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	replyID := c.Query("replyID")
	if commentID == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	rcid, err := strconv.ParseUint(replyID, 10, 64)

	// check if the comment replyes exist
	origCommentReplyes := models.Replyes{}
	err = server.DB.Debug().Model(models.Replyes{}).Where("id = ? AND comment_id = ? AND profile_id = ?", rcid, cid, profileID).Take(&origCommentReplyes).Error

	if profileID != uint32(origCommentReplyes.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// Is the authenticated user, the owner of this replye?
	if profileID != uint32(origCommentReplyes.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If all the conditions are met, delete the post
	_, err = origCommentReplyes.DeleteAReplyes(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Comment deleted",
	})
}
