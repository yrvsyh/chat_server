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

	ws := r.Group("/ws")
	ws.Use(middleware.JWTAuthMiddleware())
	ws.GET("/chat", controller.ChatHandle)

	auth := r.Group("/auth")
	auth.POST("/login", controller.Login)
	auth.POST("/register", controller.Register)

	user := r.Group("/user/:id")
	user.Use(middleware.JWTAuthMiddleware())
	user.GET("/avatar", controller.GetUserAvatar)

	return r
}
