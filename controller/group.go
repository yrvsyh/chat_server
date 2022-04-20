package controller

import (
	"chat_server/model"
	"chat_server/utils"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

func (GroupController) GetGroupAvatar(c *gin.Context) {
	form := struct {
		GroupID uint32 `form:"group_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	group, err := groupService.GetGroupByID(form.GroupID)
	if group.Avatar == "" {
		group.Avatar = "default.jpg"
	}
	filePath := "./static/avatar/group/" + group.Avatar
	if err != nil || !utils.FileExist(filePath) {
		// c.Status(http.StatusNotFound)
		// return
		filePath = "./static/avatar/group/default.jpg"
	}

	c.File(filePath)
}

func (GroupController) CreateGroupWithPublicKey(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		GroupName string `json:"group_name" binding:"required"`
		PublicKey string `json:"public_key" binding:"-"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	groupID, err := groupService.CreateGroupWithKey(json.GroupName, id, json.PublicKey)
	if err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"id": groupID})
}

func (GroupController) InvteUser(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		GroupID  uint32 `json:"group_id" binding:"required"`
		MemberID uint32 `json:"member_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	if err := groupService.InvteUser(id, json.GroupID, json.MemberID); err != nil {
		Err(c, err)
		return
	}

	Success(c)
}

func (GroupController) GetJoinedGroupInfo(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	groups, err := groupService.GetJoinedGroupsInfo(id)
	if err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"count": len(groups), "list": groups})
}

func (GroupController) UpdateGroupRemark(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		GroupID uint32 `json:"group_id" binding:"required"`
		Remark  string `json:"remark" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	if !groupService.CheckUserInGroup(json.GroupID, id) {
		Error(c, nil, "用户未加入组")
		return
	}

	userGroup, err := groupService.GetGroupUser(json.GroupID, id)
	if err != nil {
		Err(c, err)
		return
	}

	userGroup.Remark = json.Remark
	if err := groupService.UpdateGroupUser(userGroup); err != nil {
		Err(c, err)
		return
	}

	Success(c)
}

func (GroupController) CreateGroup(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		Name  string `json:"name" binding:"required"`
		Type  string `json:"type" binding:"required"`
		Label string `json:"label" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	group := &model.Group{}
	group.OwnerID = id
	group.Name = json.Name
	group.Type = json.Type
	group.Label = json.Label
	if err := groupService.CreateGroup(group); err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"id": group.ID})
}

func (GroupController) GetGroupMembers(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	form := struct {
		GroupID uint32 `form:"group_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	if !groupService.CheckUserInGroup(form.GroupID, id) {
		Error(c, nil, "用户不在此组")
		return
	}

	groupUsers, err := groupService.GetGroupMembersInfo(form.GroupID)
	if err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"count": len(groupUsers), "list": groupUsers})
}
