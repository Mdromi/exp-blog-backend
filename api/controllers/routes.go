package controllers

import (
	"github.com/Mdromi/exp-blog-backend/api/middlewares"
	_ "github.com/Mdromi/exp-blog-backend/api/middlewares"
)

func (s *Server) initializeRoutes() {
	v1 := s.Router.Group("/api/v1")
	{
		// Login Route
		v1.POST("/login", s.Login)

		// Reset Password
		v1.POST("/password/forgot", s.ForgotPassword)
		v1.POST("/password/reset", s.ResetPassword)

		// Users routes
		v1.POST("/users", s.CreateUser)
		v1.GET("/users", s.GetUsers)
		v1.GET("/users/:id", s.GetUser)
		v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		v1.PUT("/avatar/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateAvatar)
		v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		// Profiles routes
		v1.POST("/profiles", s.CreateUserProfile)
		v1.GET("/profiles", s.GetUserProfiles)
		v1.GET("/profiles/:id", s.GetUserProfile)
		v1.PUT("/profiles/:id", middlewares.TokenAuthMiddleware(), s.UpdateAUserProfile)
		v1.PUT("/avatar/profiles/:id", middlewares.TokenAuthMiddleware(), s.UpdateUserProfilePic)
		v1.DELETE("/profiles/:id", middlewares.TokenAuthMiddleware(), s.DeleteUserProfile)

		// Posts routes
		v1.POST("/posts", middlewares.TokenAuthMiddleware(), s.CreatePost)
		v1.GET("/posts", s.GetPosts)
		v1.GET("/posts/:id", s.GetPost)
		v1.PUT("/posts/:id", middlewares.TokenAuthMiddleware(), s.UpdatePost)
		v1.DELETE("/posts/:id", middlewares.TokenAuthMiddleware(), s.DeletePost)
		v1.GET("/user_posts/:id", s.GetUserProfilePosts)

		// Like Routes
		v1.GET("/likes/:id", s.GetLikes)
		v1.POST("/likes/:id", middlewares.TokenAuthMiddleware(), s.LikePost)
		v1.DELETE("/likes/:id", middlewares.TokenAuthMiddleware(), s.UnLikePost)

		// Comment routes
		v1.POST("/comments/:id", middlewares.TokenAuthMiddleware(), s.CreateComment)
		v1.GET("/comments/:id", s.GetComments)
		v1.PUT("/comments/:id/", middlewares.TokenAuthMiddleware(), s.UpdateComment)
		v1.DELETE("/comments/:id", middlewares.TokenAuthMiddleware(), s.DeleteComment)

		// Comment Replyes routes
		v1.POST("/comment/replyes/:id", middlewares.TokenAuthMiddleware(), s.CreateCommentReplye)
		v1.GET("/comments/replyes/:id", s.GetCommentReplyes)
		v1.PUT("/comments/replyes/:id/", middlewares.TokenAuthMiddleware(), s.UpdateACommentReplyes)
		v1.DELETE("/comments/replyes/:id", middlewares.TokenAuthMiddleware(), s.DeleteCommentReplye)
	}
}
