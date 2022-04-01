package controller

import (
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type GroupController struct{}

func (GroupController) CreateGroup(c *gin.Context) {
	id, _ := GetLoginUserInfo(c)

	form := struct {
		GroupName string                `form:"group_name" binding:"required"`
		PublicKey *multipart.FileHeader `form:"public_key" binding:"-"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	keyFile, err := form.PublicKey.Open()
	if err != nil {
		Err(c, err)
		return
	}

	publicKey, err := io.ReadAll(keyFile)
	if err != nil {
		Err(c, err)
		return
	}

	groupID, err := groupService.CreateGroup(form.GroupName, id, publicKey)
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

	form := struct {
		GroupID  uint32 `form:"group_id" binding:"required"`
		MemberID uint32 `form:"member_id" binding:"required"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	if err := groupService.InvteUser(id, form.GroupID, form.MemberID); err != nil {
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
