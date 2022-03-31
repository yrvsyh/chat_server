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
		Err(c, err)
		return
	}

	if file.Size > config.AvatarFileSizeLimit {
		Error(c, "file too large")
		return
	}

	path := config.AvatarPathPrefix + string(os.PathSeparator) + file.Filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		Err(c, err)
		return
	}

	if userService.UpdateUserAvatarByID(id, path) != nil {
		Err(c, err)
		return
	}

	Success(c)
}

func (UserController) GetUserFriends(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := userService.GetUserFriendsByID(id)
	if err != nil {
		Err(c, err)
		return
	}
	SuccessData(c, gin.H{"count": len(friends), "list": friends})
}

func (UserController) GetUserFriendsDetail(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := userService.GetUserFriendsDetailByID(id)
	if err != nil {
		Err(c, err)
		return
	}
	SuccessData(c, gin.H{"count": len(friends), "list": friends})
}

func (UserController) AddUserFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.PostForm("friend_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	if err := userService.AddUserFriend(id, uint32(friendID)); err != nil {
		Err(c, err)
		return
	}

	// 推送添加好友请求
	msg := &message.Message{
		Type: message.Type_FRIEND_REQUEST,
		From: id,
		To:   uint32(friendID),
	}
	chat_service.SendMessage(msg)

	Success(c)
}

func (UserController) AcceptUserFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.PostForm("friend_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	if err := userService.AcceptUserFriend(id, uint32(friendID)); err != nil {
		Err(c, err)
		return
	}

	msg := &message.Message{
		Type: message.Type_FRIEND_ACCEPT,
		From: id,
		To:   uint32(friendID),
	}
	chat_service.SendMessage(msg)

	Success(c)
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
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"remark": userFriend.Remark})
}

func (UserController) UpdateFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friendID, err := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	if err != nil {
		Err(c, err)
		return
	}

	remark := c.PostForm("remark")
	if remark == "" {
		Error(c, "invalid post data")
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, uint32(friendID))
	if err != nil {
		Err(c, err)
		return
	}

	userFriend.Remark = remark
	if err := userService.UpdateUserFriend(userFriend); err != nil {
		Err(c, err)
		return
	}

	Success(c)
}
