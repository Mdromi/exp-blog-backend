package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Mdromi/exp-blog-backend/api/auth"
	"github.com/Mdromi/exp-blog-backend/api/models"
	"github.com/Mdromi/exp-blog-backend/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

func (server *Server) LikePost(c *gin.Context) {
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
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the user exist:
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

	// Extrect Action
	// GET /like/123?action=like
	action := c.Query("action")
	if action == "" {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	like := models.LikeDislike{}
	like.ProfileID = profile.ID
	like.PostID = post.ID
	// TASK: Need to chnage
	like.Action = action

	likeCreated, err := like.SaveLike(server.DB)
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
		"response": likeCreated,
	})
}

func (server *Server) GetLikes(c *gin.Context) {
	// clear previous error if any
	errList = map[string]string{}

	postID := c.Param("id")

	// is a valid post id given to us?
	pid, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		fmt.Println("this is the error: ", err)
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// check if the post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		errList["No_post"] = "No Post Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	like := models.LikeDislike{}

	likes, err := like.GetLikesInfo(server.DB, uint(pid))
	if err != nil {
		errList["No_likes"] = "No Likes found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": likes,
	})
}

func (server *Server) UnLikePost(c *gin.Context) {
	likeID := c.Param("id")
	// is a valid like id given to us?
	lid, err := strconv.ParseUint(likeID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	// Is this user authenticated?
	profileID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	like := models.LikeDislike{}
	err = server.DB.Debug().Model(models.LikeDislike{}).Where("id = ?", lid).Take(&like).Error
	if err != nil {
		errList["No_like"] = "No Like Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	// Is the authenticated user, the owner of this post?
	if profileID != uint32(like.ProfileID) {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If all the conditions are met, delete the post
	_, err = like.DeleteLike(server.DB)
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
		"response": "Like deleted",
	})
}
