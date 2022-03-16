package controller

import (
	"chat_server/service"
	"chat_server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserAvatar(c *gin.Context) {
	username := c.Param("username")

	user, err := service.GetUserByName(username)
	if err != nil || !utils.FileExist(user.Avatar) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(user.Avatar)
}

func GetUserFriends(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := service.GetUserFriendsById(id)
	if err != nil {
		utils.Error(c, -1, "获取好友信息失败")
		return
	}
	utils.SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}

func GetUserFriendsDetail(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := service.GetUserFriendsDetailById(id)
	if err != nil {
		utils.Error(c, -1, "获取好友信息失败")
		return
	}
	utils.SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}
