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

	r.GET("/ws/chat", middleware.SessionAuthMiddleware(), controller.ChatHandle)

	auth := r.Group("/auth")
	auth.POST("/register", controller.Register)
	auth.POST("/login", controller.Login)
	auth.GET("/logout", middleware.SessionAuthMiddleware(), controller.Logout)

	user := r.Group("/user")
	user.Use(middleware.SessionAuthMiddleware())
	user.GET("/avatar/:username", controller.GetUserAvatar)
	user.POST("/avatar/upload", controller.UploadUserAvatar)
	user.GET("/friends", controller.GetUserFriends)
	user.GET("/friends_detail", controller.GetUserFriendsDetail)

	return r
}
