package router

import (
	"chat_server/controller"
	"chat_server/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	authController  = controller.AuthController{}
	userController  = controller.UserController{}
	groupController = controller.GroupController{}
	wsController    = controller.WSController{}
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(middleware.LoggerMiddleware(log.StandardLogger()))
	r.Use(gin.Recovery())

	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "./static")

	ws := r.Group("/ws")
	{
		ws.GET("/chat", middleware.SessionAuthMiddleware(), wsController.ChatHandle)
	}

	auth := r.Group("/auth")
	{
		auth.GET("/register", authController.Register)
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.GET("/logout", middleware.SessionAuthMiddleware(), authController.Logout)
	}

	user := r.Group("/user")
	{
		user.Use(middleware.SessionAuthMiddleware())
		user.GET("/search", userController.SearchUser)
		user.GET("/public_key", userController.GetUserPublicKey)
		user.GET("/avatar", userController.GetUserAvatar)
		user.POST("/avatar", userController.UploadUserAvatar)
		user.GET("/friends", userController.GetUserFriends)
		user.GET("/friends_detail", userController.GetUserFriendsDetail)
		user.POST("/add_friend", userController.AddUserFriend)
		user.POST("/accept_friend", userController.AcceptUserFriend)
		user.GET("/remark", userController.GetFriendRemark)
		user.POST("/remark", userController.UpdateFriendRemark)
	}

	group := r.Group("/group")
	{
		group.Use(middleware.SessionAuthMiddleware())
		group.POST("/avatar", groupController.GetGroupAvatar)
		group.POST("/invite", groupController.InvteUser)
		group.POST("/create", groupController.CreateGroup)
		group.GET("/joined_group", groupController.GetJoinedGroupInfo)
	}

	return r
}
