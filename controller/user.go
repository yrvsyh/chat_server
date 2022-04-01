package controller

import (
	"chat_server/config"
	"chat_server/message"
	"chat_server/utils"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (UserController) GetUserAvatar(c *gin.Context) {
	form := struct {
		UserID uint32 `form:"user_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(form.UserID)
	if err != nil || !utils.FileExist(user.Avatar) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(user.Avatar)
}

func (UserController) UploadUserAvatar(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	form := struct {
		Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	file := form.Avatar
	if file.Size > config.AvatarFileSizeLimit {
		Error(c, "file too large")
		return
	}

	path := config.AvatarPathPrefix + string(os.PathSeparator) + file.Filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		Err(c, err)
		return
	}

	if err := userService.UpdateUserAvatarByID(id, path); err != nil {
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

	form := struct {
		FriendID uint32 `form:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	if err := userService.AddUserFriend(id, form.FriendID); err != nil {
		Err(c, err)
		return
	}

	// 推送添加好友请求
	msg := &message.Message{
		Id:   time.Now().UnixNano(),
		Type: message.Type_FRIEND_REQUEST,
		From: id,
		To:   form.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

func (UserController) AcceptUserFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	form := struct {
		FriendID uint32 `form:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	if err := userService.AcceptUserFriend(id, form.FriendID); err != nil {
		Err(c, err)
		return
	}

	msg := &message.Message{
		Id:   time.Now().UnixNano(),
		Type: message.Type_FRIEND_ACCEPT,
		From: id,
		To:   form.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

func (UserController) GetFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	form := struct {
		FriendID uint32 `form:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, form.FriendID)
	if err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"remark": userFriend.Remark})
}

func (UserController) UpdateFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	form := struct {
		FriendID uint32 `form:"friend_id" binding:"required"`
		Remark   string `form:"remark" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, form.FriendID)
	if err != nil {
		Err(c, err)
		return
	}

	userFriend.Remark = form.Remark
	if err := userService.UpdateUserFriend(userFriend); err != nil {
		Err(c, err)
		return
	}

	Success(c)
}
