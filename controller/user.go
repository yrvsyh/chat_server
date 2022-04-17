package controller

import (
	"chat_server/config"
	"chat_server/message"
	"chat_server/utils"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

func (UserController) SearchUser(c *gin.Context) {
	query := struct {
		Name string `form:"name" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		Err(c, err)
		return
	}

	users, err := userService.SearchUserByName(query.Name)
	if err != nil {
		Err(c, err)
	}

	SuccessData(c, gin.H{"count": len(users), "list": users})
}

func (UserController) GetUserPublicKey(c *gin.Context) {
	query := struct {
		ID uint32 `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(query.ID)
	if err != nil {
		Err(c, err)
	}

	SuccessData(c, gin.H{"public_key": user.PublicKey})
}

func (UserController) GetUserAvatar(c *gin.Context) {
	query := struct {
		ID uint32 `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBindQuery(&query); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(query.ID)
	if user.Avatar == "" {
		user.Avatar = "default.jpg"
	}
	filePath := "./static/avatar/" + user.Avatar
	if err != nil || !utils.FileExist(filePath) {
		// c.Status(http.StatusNotFound)
		// return
		filePath = "./static/avatar/default.jpg"
	}

	c.File(filePath)
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

	json := struct {
		FriendID uint32 `json:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	if err := userService.AddUserFriend(id, json.FriendID); err != nil {
		Err(c, err)
		return
	}

	// 推送添加好友请求
	msg := &message.Message{
		Type: message.Type_FRIEND_REQUEST,
		From: id,
		To:   json.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

func (UserController) AcceptUserFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		FriendID uint32 `json:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	if err := userService.AcceptUserFriend(id, json.FriendID); err != nil {
		Err(c, err)
		return
	}

	msg := &message.Message{
		Type: message.Type_FRIEND_ACCEPT,
		From: id,
		To:   json.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

func (UserController) GetFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		FriendID uint32 `json:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, json.FriendID)
	if err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"remark": userFriend.Remark})
}

func (UserController) UpdateFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		FriendID uint32 `json:"friend_id" binding:"required"`
		Remark   string `json:"remark" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, json.FriendID)
	if err != nil {
		Err(c, err)
		return
	}

	userFriend.Remark = json.Remark
	if err := userService.UpdateUserFriend(userFriend); err != nil {
		Err(c, err)
		return
	}

	Success(c)
}
