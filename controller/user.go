package controller

import (
	"chat_server/config"
	"chat_server/service"
	"chat_server/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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

func UploadUserAvatar(c *gin.Context) {
	username := GetLoginUserName(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		utils.Error(c, -1, "获取文件失败")
		return
	}

	if file.Size > config.AvatarFileSizeLimit {
		utils.Error(c, -1, "文件尺寸过大")
		return
	}

	path := config.AvatarPathPrefix + string(os.PathSeparator) + file.Filename

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		utils.Error(c, -1, "保存头像失败")
		return
	}

	if service.UpdateUserAvatarByName(username, path) != nil {
		utils.Error(c, -1, "用户头像更新失败")
		return
	}

	utils.Success(c, "上传头像成功")
}

func GetUserFriends(c *gin.Context) {
	username := GetLoginUserName(c)

	friends, err := service.GetUserFriendsByName(username)
	if err != nil {
		utils.Error(c, -1, "获取好友信息失败")
		return
	}
	utils.SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}

func GetUserFriendsDetail(c *gin.Context) {
	username := GetLoginUserName(c)

	friends, err := service.GetUserFriendsDetailByName(username)
	if err != nil {
		utils.Error(c, -1, "获取好友信息失败")
		return
	}
	utils.SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}
