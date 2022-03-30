package router

import (
	"chat_server/controller"
	"chat_server/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware(log.StandardLogger()))
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/**/*")

	r.GET("/ws/chat", middleware.SessionAuthMiddleware(), controller.ChatController{}.ChatHandle)

	auth := r.Group("/auth")
	{
		auth.GET("/register", controller.AuthController{}.Register)
		auth.POST("/register", controller.AuthController{}.Register)
		auth.POST("/login", controller.AuthController{}.Login)
		auth.GET("/logout", middleware.SessionAuthMiddleware(), controller.AuthController{}.Logout)
	}

	user := r.Group("/user")
	{
		user.Use(middleware.SessionAuthMiddleware())
		user.GET("/avatar/:user_id", controller.UserController{}.GetUserAvatar)
		user.POST("/avatar/upload", controller.UserController{}.UploadUserAvatar)
		user.GET("/friends", controller.UserController{}.GetUserFriends)
		user.GET("/friends_detail", controller.UserController{}.GetUserFriendsDetail)
		user.POST("/add_friends", controller.UserController{}.AddUserFriend)
		user.POST("/accept_friends", controller.UserController{}.AcceptUserFriend)
		user.GET("/remark/:friend_id", controller.UserController{}.GetFriendRemark)
		user.POST("/remark/:friend_id", controller.UserController{}.UpdateFriendRemark)
	}

	return r
}
