package controller

import (
	"chat_server/config"
	"chat_server/message"
	"chat_server/model"
	"chat_server/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct{}

func (UserController) SearchUser(c *gin.Context) {
	form := struct {
		Name string `form:"name" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	users, err := userService.SearchUserByName(form.Name)
	if err != nil {
		Err(c, err)
	}

	SuccessData(c, gin.H{"count": len(users), "list": users})
}

func (UserController) GetUserInfo(c *gin.Context) {
	form := struct {
		ID uint32 `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(form.ID)
	if err != nil {
		Err(c, err)
	}

	ret := gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
	}
	SuccessData(c, ret)
}

func (UserController) UpdateUserInfo(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	user := &model.User{}
	user.ID = id
	user.Nickname = json.Nickname
	user.Email = json.Email
	user.Phone = json.Phone

	if err := userService.UpdateUser(user); err != nil {
		Err(c, err)
	}

	Success(c)
}

func (UserController) GetUserPublicKey(c *gin.Context) {
	form := struct {
		ID uint32 `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(form.ID)
	if err != nil {
		Err(c, err)
	}

	SuccessData(c, gin.H{"public_key": user.PublicKey})
}

func (UserController) GetUserAvatar(c *gin.Context) {
	form := struct {
		ID uint32 `form:"id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.GetUserByID(form.ID)
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

	// form := struct {
	// 	Avatar *multipart.FileHeader `form:"avatar" binding:"required"`
	// }{}

	// if err := c.ShouldBind(&form); err != nil {
	// 	Err(c, err)
	// 	return
	// }

	file, err := c.FormFile("avatar")
	if err != nil {
		Err(c, err)
		return
	}
	logrus.Info(file.Filename)

	// file := form.Avatar
	if file.Size > config.AvatarFileSizeLimit {
		Error(c, "file too large")
		return
	}

	var fileName string
	var path string
	for {
		// fileName = utils.GenUUID() + "-" + file.Filename
		fileName = utils.GenUUID()
		path = config.AvatarPathPrefix + string(os.PathSeparator) + fileName
		fileInfo, err := os.Stat(path)
		if err != nil || fileInfo.IsDir() {
			break
		}
	}

	if err := c.SaveUploadedFile(file, path); err != nil {
		Err(c, err)
		return
	}

	if err := userService.UpdateUserAvatarByID(id, fileName); err != nil {
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
		// Err(c, err)
		// return
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
