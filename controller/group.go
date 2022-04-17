package controller

import (
	"github.com/gin-gonic/gin"
)

type GroupController struct{}

func (GroupController) CreateGroup(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	json := struct {
		GroupName string `json:"group_name" binding:"required"`
		PublicKey string `json:"public_key" binding:"-"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	groupID, err := groupService.CreateGroup(json.GroupName, id, json.PublicKey)
	if err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"id": groupID})
}

func (GroupController) GetGroupAvatar(c *gin.Context) {
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
