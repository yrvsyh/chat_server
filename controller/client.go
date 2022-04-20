package controller

import (
	"chat_server/config"
	"chat_server/message"
	"chat_server/model"
	"chat_server/utils"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type ClientController struct{}

// 获取用户的头像, 返回头像文件
func (ClientController) GetUserAvatar(c *gin.Context) {
	data := struct {
		UserID uint32 `form:"user_id" json:"user_id" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	user, err := userService.GetUserByID(data.UserID)
	if err != nil {
		Error(c, err, "获取用户失败")
		return
	}

	fileName := user.Avatar

	if strings.TrimSpace(fileName) == "" {
		fileName = "default.jpg"
	}

	filePath := config.AvatarPathPrefix + fileName
	if !utils.FileExist(filePath) {
		Error(c, nil, "头像不存在")
		return
	}

	c.File(filePath)
}

// 获取组的头像, 返回头像文件
func (ClientController) GetGroupAvatar(c *gin.Context) {
	data := struct {
		GroupID uint32 `form:"group_id" json:"group_id" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	user, err := groupService.GetGroupByID(data.GroupID)
	if err != nil {
		Error(c, err, "获取群组信息失败")
		return
	}

	fileName := user.Avatar

	if strings.TrimSpace(fileName) == "" {
		fileName = "default.jpg"
	}

	filePath := config.AvatarPathPrefix + fileName
	if !utils.FileExist(filePath) {
		Error(c, nil, "头像不存在")
		return
	}

	c.File(filePath)
}

// 返回好友的基本信息 id name(username nickname(对方的昵称) remark(对好友的备注)) avatar public_key
func (ClientController) GetUserFriends(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	friends, err := userService.GetUserFriendsByID(id)
	if err != nil {
		Error(c, err, "获取好友信息失败")
		return
	}

	type VO struct {
		ID        uint32 `json:"id"`
		Name      string `json:"name"`
		Avatar    string `json:"avatar"`
		PublicKey string `json:"public_key"`
	}
	var ret []VO

	for _, friend := range friends {
		vo := VO{
			ID:        friend.UserID,
			Name:      friend.User.Username,
			PublicKey: friend.User.PublicKey,
		}
		if name := strings.TrimSpace(friend.User.Nickname); name != "" {
			vo.Name = name
		}
		if name := strings.TrimSpace(friend.Remark); name != "" {
			vo.Name = name
		}
		ret = append(ret, vo)
	}

	SuccessData(c, gin.H{"count": len(ret), "list": ret})
}

// 返回加入的组的基本信息 id name(groupname remark(对组的备注)) nickname(用户在组内的昵称) avatar
func (ClientController) GetUserGroups(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	groups, err := groupService.GetJoinedGroupsInfo(id)
	if err != nil {
		Error(c, err, "获取群组信息失败")
		return
	}

	type VO struct {
		ID       uint32 `json:"id"`
		Name     string `json:"name"`
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
	}
	var ret []VO

	for _, group := range groups {
		vo := VO{
			ID:       group.GroupID,
			Name:     group.Group.Name,
			Avatar:   group.Group.Avatar,
			Nickname: group.Nickname,
		}
		if name := strings.TrimSpace(group.Remark); name != "" {
			vo.Name = name
		}
		ret = append(ret, vo)
	}

	SuccessData(c, gin.H{"count": len(ret), "list": ret})
}

// 获取组内所有成员信息 id name(username nickname nickname_in_group) avatar
// 收到消息后通过GetGroupMemberInfo加载
func (ClientController) GetGroupMembers(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		GroupID uint32 `form:"group_id" json:"group_id" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if !groupService.CheckUserInGroup(data.GroupID, id) {
		Error(c, nil, "用户不在此组")
	}

	groupUsers, err := groupService.GetGroupMembersInfo(data.GroupID)
	if err != nil {
		Error(c, err, "获取组成员失败")
		return
	}

	type VO struct {
		ID     uint32 `json:"id"`
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}
	var ret []VO

	for _, groupUser := range groupUsers {
		vo := VO{
			ID:     groupUser.UserID,
			Name:   groupUser.User.Username,
			Avatar: groupUser.User.Avatar,
		}
		if name := strings.TrimSpace(groupUser.User.Nickname); name != "" {
			vo.Name = name
		}
		if name := strings.TrimSpace(groupUser.Nickname); name != "" {
			vo.Name = name
		}
		ret = append(ret, vo)
	}

	SuccessData(c, gin.H{"count": len(ret), "list": ret})
}

// 返回好友的详细信息 id username nickname ...
func (ClientController) GetFriendInfo(c *gin.Context) {
	data := struct {
		FriendID uint32 `form:"friend_id" json:"friend_id" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	// TODO 判断用户是否为当前用户好友

	user, err := userService.GetUserByID(data.FriendID)
	if err != nil {
		Error(c, err, "获取好友信息失败")
		return
	}

	type VO struct {
		ID       uint32 `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
	}
	var ret = VO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Email:    user.Email,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	}

	SuccessData(c, ret)
}

// 返回与当前用户在同一组的用户详细信息 id username nickname nickname_in_group
func (ClientController) GetGroupMemberInfo(c *gin.Context) {
	data := struct {
		MemberID uint32 `form:"member_id" json:"member_id" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	// TODO 判断用户是否与当前用户存在共同小组

	user, err := userService.GetUserByID(data.MemberID)
	if err != nil {
		Error(c, err, "获取用户信息失败")
		return
	}

	type VO struct {
		ID       uint32 `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	var ret = VO{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}

	SuccessData(c, ret)
}

// 发送图片到指定用户, 服务端构造推送信息
func (ClientController) SendImage(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		ToFriend bool                  `form:"to_friend" json:"to_friend" binding:"required"`
		To       uint32                `form:"to" json:"to" binding:"required"`
		Image    *multipart.FileHeader `form:"image" json:"image" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if data.ToFriend {
		// 判断是否为好友
		if _, err := userService.GetUserFriendDetailByFriendID(id, data.To); err != nil {
			Error(c, err, "获取用户好友失败")
			return
		}
	} else {
		// 判断是否在组内
		if !groupService.CheckUserInGroup(data.To, id) {
			Error(c, nil, "用户不在组内")
			return
		}
	}

	file := data.Image

	var fileName string
	var filePath string
	for {
		fileName = utils.GenUUID()
		filePath = config.ImagePathPrefix + fileName
		if !utils.FileExist(filePath) {
			break
		}
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Error(c, err, "保存图片失败")
		return
	}

	// if data.ToFriend {
	// 	msg := &message.Message{
	// 		Type:    message.Type_FRIEND_IMAGE,
	// 		From:    id,
	// 		To:      data.To,
	// 		Content: []byte(fileName),
	// 	}
	// 	manager.SendMessage(msg)
	// } else {
	// 	msg := &message.Message{
	// 		Type:    message.Type_GROUP_IMAGE,
	// 		From:    id,
	// 		To:      data.To,
	// 		Content: []byte(fileName),
	// 	}
	// 	manager.SendMessage(msg)
	// }

	SuccessData(c, gin.H{"filename": fileName})
}

func (ClientController) SendFile(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		ToFriend bool                  `form:"to_friend" json:"to_friend" binding:"required"`
		To       uint32                `form:"to" json:"to" binding:"required"`
		File     *multipart.FileHeader `form:"file" json:"file" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if data.ToFriend {
		// 判断是否为好友
		if _, err := userService.GetUserFriendDetailByFriendID(id, data.To); err != nil {
			Error(c, err, "获取用户好友失败")
			return
		}
	} else {
		// 判断是否在组内
		if !groupService.CheckUserInGroup(data.To, id) {
			Error(c, nil, "用户不在组内")
			return
		}
	}

	file := data.File

	var fileName string
	var filePath string
	for {
		fileName = utils.GenUUID()
		filePath = config.ImagePathPrefix + fileName
		if !utils.FileExist(filePath) {
			break
		}
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Error(c, err, "保存文件失败")
		return
	}

	// if data.ToFriend {
	// 	msg := &message.Message{
	// 		Type:    message.Type_FRIEND_FILE,
	// 		From:    id,
	// 		To:      data.To,
	// 		Content: []byte(fileName),
	// 	}
	// 	manager.SendMessage(msg)
	// } else {
	// 	msg := &message.Message{
	// 		Type:    message.Type_GROUP_FILE,
	// 		From:    id,
	// 		To:      data.To,
	// 		Content: []byte(fileName),
	// 	}
	// 	manager.SendMessage(msg)
	// }

	SuccessData(c, gin.H{"filename": fileName})
}

// 更新当前用户信息
func (ClientController) UpdateUserInfo(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		Nickname string `form:"nickname" json:"nickname" binding:"required"`
		Email    string `form:"email" json:"email" binding:"required"`
		Phone    string `form:"phone" json:"phone" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	user, err := userService.GetUserByID(id)
	if err != nil {
		Error(c, err, "获取用户信息失败")
		return
	}

	user.Nickname = data.Nickname
	user.Email = data.Email
	user.Phone = data.Phone
	if err := userService.UpdateUser(user); err != nil {
		Error(c, err, "更新用户信息失败")
		return
	}

	Success(c)
}

// 更新当前用户头像
func (ClientController) UpdateUserAvatar(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		Avatar *multipart.FileHeader `form:"avatar" json:"avatar" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	file := data.Avatar
	if file.Size > config.AvatarFileSizeLimit {
		Error(c, nil, "头像文件过大")
		return
	}

	user, err := userService.GetUserByID(id)
	if err != nil {
		Error(c, err, "获取用户信息失败")
		return
	}

	// 删除旧文件
	os.Remove(config.AvatarPathPrefix + user.Avatar)

	var fileName string
	var filePath string
	for {
		fileName = utils.GenUUID()
		filePath = config.AvatarPathPrefix + fileName
		if !utils.FileExist(filePath) {
			break
		}
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Error(c, err, "保存头像文件失败")
		return
	}

	user.Avatar = fileName
	if err := userService.UpdateUser(user); err != nil {
		Error(c, err, "更新用户头像失败")
		return
	}

	Success(c)
}

// 修改好友备注
func (ClientController) UpdateFriendRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		FriendID uint32 `form:"friend_id" json:"friend_id" binding:"required"`
		Remark   string `form:"remark" json:"remark" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Error(c, err, "参数错误")
		return
	}

	userFriend, err := userService.GetUserFriendDetailByFriendID(id, json.FriendID)
	if err != nil {
		Error(c, err, "获取好友信息失败")
		return
	}

	userFriend.Remark = json.Remark
	if err := userService.UpdateUserFriend(userFriend); err != nil {
		Error(c, err, "更新好友信息失败")
		return
	}

	Success(c)
}

// 更新群头像
func (ClientController) UpdateGroupAvatar(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		GroupID uint32                `form:"group_id" json:"group_id" binding:"required"`
		Avatar  *multipart.FileHeader `form:"avatar" json:"avatar" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	file := data.Avatar
	if file.Size > config.AvatarFileSizeLimit {
		Error(c, nil, "头像文件过大")
		return
	}

	groupUser, err := groupService.GetGroupUser(data.GroupID, id)
	if err != nil {
		Error(c, err, "获取组成员信息失败")
		return
	}

	if groupUser.Group.OwnerID != id {
		Error(c, err, "无权限")
		return
	}

	group, err := groupService.GetGroupByID(data.GroupID)
	if err != nil {
		Error(c, err, "获取群组信息失败")
		return
	}

	// 删除旧文件
	os.Remove(config.AvatarPathPrefix + group.Avatar)

	var fileName string
	var filePath string
	for {
		fileName = utils.GenUUID()
		filePath = config.AvatarPathPrefix + fileName
		if !utils.FileExist(filePath) {
			break
		}
	}

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		Error(c, err, "保存头像文件失败")
		return
	}

	group.Avatar = fileName
	if err := groupService.UpdateGroup(group); err != nil {
		Error(c, err, "更新群组头像失败")
		return
	}

	Success(c)
}

// 修改群备注
func (ClientController) UpdateGroupRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		GroupID uint32 `form:"group_id" json:"group_id" binding:"required"`
		Remark  string `form:"remark" json:"remark" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	groupUser, err := groupService.GetGroupUser(data.GroupID, id)
	if err != nil {
		Error(c, err, "获取组成员信息失败")
		return
	}

	groupUser.Remark = data.Remark
	if err := groupService.UpdateGroupUser(groupUser); err != nil {
		Error(c, err, "更新组成员信息失败")
		return
	}

	Success(c)
}

// 修改群内昵称
func (ClientController) UpdateNicknameInGroup(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		GroupID  uint32 `form:"group_id" json:"group_id" binding:"required"`
		Nickname string `form:"nickname" json:"nickname" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	groupUser, err := groupService.GetGroupUser(data.GroupID, id)
	if err != nil {
		Error(c, err, "获取组成员信息失败")
		return
	}

	groupUser.Nickname = data.Nickname
	if err := groupService.UpdateGroupUser(groupUser); err != nil {
		Error(c, err, "更新组成员信息失败")
		return
	}

	Success(c)
}

// 搜索用户
func (ClientController) SearchUser(c *gin.Context) {
	data := struct {
		Name string `form:"name" json:"name" binding:"required"`
	}{}
	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	users, err := userService.SearchUserByName(data.Name)
	if err != nil {
		Error(c, err, "搜索用户失败")
		return
	}

	type VO struct {
		ID       uint32 `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
	}
	var ret []VO

	for _, user := range users {
		vo := VO{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
		}
		ret = append(ret, vo)
	}

	SuccessData(c, gin.H{"count": len(ret), "list": ret})
}

// 添加好友
func (ClientController) AddFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		FriendID uint32 `form:"friend_id" json:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if err := userService.AddUserFriend(id, data.FriendID); err != nil {
		Error(c, err, "添加好友失败")
		return
	}

	// 推送添加好友请求
	msg := &message.Message{
		Type: message.Type_FRIEND_REQUEST,
		From: id,
		To:   data.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

// 同意好友申请
func (ClientController) AcceptFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		FriendID uint32 `form:"friend_id" json:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if err := userService.AcceptUserFriend(id, data.FriendID); err != nil {
		Error(c, err, "接受好友请求失败")
		return
	}

	msg := &message.Message{
		Type: message.Type_FRIEND_ACCEPT,
		From: id,
		To:   data.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

// 删除好友
func (ClientController) DeleteFriend(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		FriendID uint32 `form:"friend_id" json:"friend_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if err := userService.DeleteUserFriend(id, data.FriendID); err != nil {
		Error(c, err, "删除好友失败")
		return
	}

	msg := &message.Message{
		Type: message.Type_FRIEND_DISBAND,
		From: id,
		To:   data.FriendID,
	}
	manager.SendMessage(msg)

	Success(c)
}

// 创建群聊
func (ClientController) CreateGroup(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		Name  string `form:"name" json:"name" binding:"required"`
		Type  string `form:"type" json:"type" binding:"required"`
		Label string `form:"label" json:"label" binding:"required"`
	}{}

	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	group := &model.Group{}
	group.OwnerID = id
	group.Name = data.Name
	group.Type = data.Type
	group.Label = data.Label
	if err := groupService.CreateGroup(group); err != nil {
		Error(c, err, "创建群组失败")
		return
	}

	SuccessData(c, gin.H{"id": group.ID})
}

// 邀请好友进群
func (ClientController) InviteUser(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	data := struct {
		GroupID  uint32 `form:"group_id" json:"group_id" binding:"required"`
		MemberID uint32 `form:"member_id" json:"member_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&data); err != nil {
		Error(c, err, "参数错误")
		return
	}

	if err := groupService.InvteUser(id, data.GroupID, data.MemberID); err != nil {
		Error(c, err, "邀请好友失败")
		return
	}

	Success(c)
}
