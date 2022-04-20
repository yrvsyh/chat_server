package router

import (
	"chat_server/controller"
	"chat_server/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	authController   = controller.AuthController{}
	userController   = controller.UserController{}
	groupController  = controller.GroupController{}
	clientController = controller.ClientController{}
	wsController     = controller.WSController{}
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

	client := r.Group("/client")
	{
		client.Use(middleware.SessionAuthMiddleware())

		client.GET("/user_avatar", clientController.GetUserAvatar)
		client.GET("/group_avatar", clientController.GetGroupAvatar)
		client.GET("/user_friends", clientController.GetUserFriends)
		client.GET("/user_groups", clientController.GetUserGroups)
		client.GET("/friend_info", clientController.GetFriendInfo)
		client.GET("/group_member_info", clientController.GetGroupMemberInfo)
		client.POST("/send_image", clientController.SendImage)
		client.POST("/send_file", clientController.SendFile)
		client.PUT("/user_info", clientController.UpdateUserInfo)
		client.PUT("/user_avatar", clientController.UpdateUserAvatar)
		client.PUT("/friend_remark", clientController.UpdateFriendRemark)
		client.PUT("/group_avatar", clientController.UpdateGroupAvatar)
		client.PUT("/group_remark", clientController.UpdateGroupRemark)
		client.PUT("/nickname_in_group", clientController.UpdateNicknameInGroup)
		client.GET("/search_user", clientController.SearchUser)
		client.POST("/add_friend", clientController.AddFriend)
		client.POST("/accept_friend", clientController.AcceptFriend)
		client.POST("/delete_friend", clientController.DeleteFriend)
		client.POST("/create_group", clientController.CreateGroup)
		client.POST("/invite_user", clientController.InviteUser)
	}

	// user := r.Group("/user")
	// {
	// 	user.Use(middleware.SessionAuthMiddleware())
	// 	user.GET("/info", userController.GetUserInfo)
	// 	user.POST("/update_info", userController.UpdateUserInfo)
	// 	user.GET("/search", userController.SearchUser)
	// 	user.GET("/public_key", userController.GetUserPublicKey)
	// 	user.GET("/avatar", userController.GetUserAvatar)
	// 	user.POST("/avatar", userController.UploadUserAvatar)
	// 	user.GET("/friends", userController.GetUserFriends)
	// 	user.GET("/friends_detail", userController.GetUserFriendsDetail)
	// 	user.POST("/add_friend", userController.AddUserFriend)
	// 	user.POST("/accept_friend", userController.AcceptUserFriend)
	// 	user.GET("/remark", userController.GetFriendRemark)
	// 	user.POST("/remark", userController.UpdateFriendRemark)
	// }

	// group := r.Group("/group")
	// {
	// 	group.Use(middleware.SessionAuthMiddleware())
	// 	group.GET("/avatar", groupController.GetGroupAvatar)
	// 	group.POST("/invite", groupController.InvteUser)
	// 	group.POST("/create_with_key", groupController.CreateGroupWithPublicKey)
	// 	group.GET("/joined_group", groupController.GetJoinedGroupInfo)
	// 	group.POST("/remark", groupController.UpdateGroupRemark)
	// 	group.POST("/create", groupController.CreateGroup)
	// 	group.GET("/members", groupController.GetGroupMembers)
	// }

	return r
}
