package controller

import (
	"chat_server/config"
	"chat_server/controller/chat"
	"chat_server/message"
	"chat_server/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (UserController) GetUserAvatar(c *gin.Context) {
	username := c.Param("username")

	user, err := userService.GetUserByName(username)
	if err != nil || !utils.FileExist(user.Avatar) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(user.Avatar)
}

func (UserController) UploadUserAvatar(c *gin.Context) {
	username := GetLoginUserName(c)

	file, err := c.FormFile("avatar")
	if err != nil {
		Error(c, -1, "获取文件失败")
		return
	}

	if file.Size > config.AvatarFileSizeLimit {
		Error(c, -1, "文件尺寸过大")
		return
	}

	path := config.AvatarPathPrefix + string(os.PathSeparator) + file.Filename

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		Error(c, -1, "保存头像失败")
		return
	}

	if userService.UpdateUserAvatarByName(username, path) != nil {
		Error(c, -1, "用户头像更新失败")
		return
	}

	Success(c, "上传头像成功")
}

func (UserController) GetUserFriends(c *gin.Context) {
	username := GetLoginUserName(c)

	friends, err := userService.GetUserFriendsByName(username)
	if err != nil {
		Error(c, -1, "获取好友信息失败")
		return
	}
	SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}

func (UserController) GetUserFriendsDetail(c *gin.Context) {
	username := GetLoginUserName(c)

	friends, err := userService.GetUserFriendsDetailByName(username)
	if err != nil {
		Error(c, -1, "获取好友信息失败")
		return
	}
	SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}

func (UserController) AddUserFriend(c *gin.Context) {
	username := GetLoginUserName(c)

	friendName := c.PostForm("friend_name")
	if friendName == "" {
		Error(c, -1, "参数错误")
		return
	}
	err := userService.AddUserFriend(username, friendName)
	if err != nil {
		Error(c, -1, "添加好友失败")
	}

	// 推送添加好友请求
	msg := &message.Message{
		Type: message.MessageType_FRIEND_REQUEST,
		From: username,
		To:   friendName,
	}
	chat.SendMessage(msg)

	Success(c, "成功发起请求")
}

func (UserController) AcceptUserFriend(c *gin.Context) {
	username := GetLoginUserName(c)

	friendName := c.PostForm("friend_name")
	if friendName == "" {
		Error(c, -1, "参数错误")
		return
	}

	err := userService.AcceptUserFriend(username, friendName)
	if err != nil {
		Error(c, -1, "添加好友失败")
		return
	}

	msg := &message.Message{
		Type: message.MessageType_FRIEND_ACCEPT,
		From: username,
		To:   friendName,
	}
	chat.SendMessage(msg)

	Success(c, "添加好友成功")
}

func (UserController) GetFriendRemark(c *gin.Context) {
	username := GetLoginUserName(c)
	friendName := c.Param("friend_name")
	if friendName == "" {
		Error(c, -1, "参数错误")
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendName(username, friendName)
	if err != nil {
		Error(c, -1, "查询好友失败")
		return
	}

	SuccessWithData(c, "查询备注成功", gin.H{"remark": userFriend.Remark})
}

func (UserController) UpdateFriendRemark(c *gin.Context) {
	username := GetLoginUserName(c)
	friendName := c.Param("friend_name")
	remark := c.PostForm("remark")
	if friendName == "" || remark == "" {
		Error(c, -1, "参数错误")
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendName(username, friendName)
	if err != nil {
		Error(c, -1, "查询好友失败")
		return
	}
	userFriend.Remark = remark
	err = userService.UpdateUserFriend(userFriend)
	if err != nil {
		Error(c, -1, "更新备注失败")
		return
	}

	Success(c, "更新备注成功")
}
