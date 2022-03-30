package controller

import (
	"chat_server/config"
	"chat_server/message"
	"chat_server/service/chat_service"
	"chat_server/utils"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (UserController) GetUserAvatar(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(uint32(id))
	if err != nil || !utils.FileExist(user.Avatar) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(user.Avatar)
}

func (UserController) UploadUserAvatar(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

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

	if userService.UpdateUserAvatarByID(id, path) != nil {
		Error(c, -1, "用户头像更新失败")
		return
	}

	Success(c, "上传头像成功")
}

func (UserController) GetUserFriends(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := userService.GetUserFriendsByID(id)
	if err != nil {
		Error(c, -1, "获取好友信息失败")
		return
	}
	SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}

func (UserController) GetUserFriendsDetail(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := userService.GetUserFriendsDetailByID(id)
	if err != nil {
		Error(c, -1, "获取好友信息失败")
		return
	}
	SuccessWithData(c, "获取好友信息成功", gin.H{"count": len(friends), "list": friends})
}

func (UserController) AddUserFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.PostForm("friend_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	if err := userService.AddUserFriend(id, uint32(friendID)); err != nil {
		Error(c, -1, "添加好友失败")
		return
	}

	// 推送添加好友请求
	msg := &message.Message{
		Type: message.Type_FRIEND_REQUEST,
		From: id,
		To:   uint32(friendID),
	}
	chat_service.SendMessage(msg)

	Success(c, "成功发起请求")
}

func (UserController) AcceptUserFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.PostForm("friend_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	if err := userService.AcceptUserFriend(id, uint32(friendID)); err != nil {
		Error(c, -1, "添加好友失败")
		return
	}

	msg := &message.Message{
		Type: message.Type_FRIEND_ACCEPT,
		From: id,
		To:   uint32(friendID),
	}
	chat_service.SendMessage(msg)

	Success(c, "添加好友成功")
}

func (UserController) GetFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, uint32(friendID))
	if err != nil {
		Error(c, -1, "查询好友失败")
		return
	}

	SuccessWithData(c, "查询备注成功", gin.H{"remark": userFriend.Remark})
}

func (UserController) UpdateFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	remark := c.PostForm("remark")
	if err != nil || remark == "" {
		Error(c, -1, "参数错误")
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, uint32(friendID))
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
