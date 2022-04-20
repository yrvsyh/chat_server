package service

import (
	"chat_server/model"
	"errors"

	"gorm.io/gorm"
)

type GroupService struct{}

func (GroupService) GetGroupByID(id uint32) (*model.Group, error) {
	group := &model.Group{}
	err := db.First(group, id).Error
	return group, err
}

func (GroupService) GetGroupUser(groupID uint32, userID uint32) (*model.GroupUser, error) {
	groupUser := &model.GroupUser{}
	err := db.Where("group_id = ? and user_id = ?", groupID, userID).First(groupUser).Error
	return groupUser, err
}

func (GroupService) UpdateGroupUser(groupUser *model.GroupUser) error {
	return db.Save(groupUser).Error
}

func (GroupService) GetGroupMemberSet(id uint32) (map[uint32]struct{}, error) {
	ret := make(map[uint32]struct{})
	var groupUsers []model.GroupUser
	err := db.Where("group_id = ?", id).Find(&groupUsers).Error
	for _, groupUser := range groupUsers {
		ret[groupUser.UserID] = struct{}{}
	}
	return ret, err
}

func (GroupService) CheckUserInGroup(groupID uint32, userID uint32) bool {
	groupUser := &model.GroupUser{}
	err := db.Where("group_id = ? and user_id = ?", groupID, userID).First(groupUser).Error
	return err == nil
}

func (GroupService) CreateGroupWithKey(name string, id uint32, publicKey string) (uint32, error) {
	group := &model.Group{Name: name, OwnerID: id, PublicKey: publicKey}
	err := db.Create(group).Error
	return group.ID, err
}

func (GroupService) InvteUser(initiatorID uint32, groupID uint32, inviteesID uint32) error {
	initiator, err := userService.GetUserByID(initiatorID)
	if err != nil {
		return err
	}

	group, err := groupService.GetGroupByID(groupID)
	if err != nil {
		return err
	}

	if initiator.ID != group.OwnerID && !groupService.CheckUserInGroup(group.ID, initiator.ID) {
		return errors.New("permission denied")
	}

	invitees, err := userService.GetUserByID(inviteesID)
	if err != nil {
		return err
	}

	groupUser := &model.GroupUser{
		GroupID: groupID,
		UserID:  invitees.ID,
	}
	if err := db.Create(groupUser).Error; err != nil {
		return err
	}

	if initiatorID != group.OwnerID && groupService.CheckUserInGroup(groupID, initiatorID) {
		// TODO 管理员审核
	}

	return nil
}

func (GroupService) GetJoinedGroupsInfo(id uint32) ([]model.GroupUser, error) {
	var groupUsers []model.GroupUser
	err := db.Preload("Group").Where("user_id = ?", id).Find(&groupUsers).Error
	return groupUsers, err
}

func (GroupService) CreateGroup(group *model.Group) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		groupUser := &model.GroupUser{}
		groupUser.GroupID = group.ID
		groupUser.UserID = group.OwnerID
		if err := tx.Create(groupUser).Error; err != nil {
			return err
		}

		return nil
	})

}

func (GroupService) GetGroupMembersInfo(groupID uint32) ([]model.GroupUser, error) {
	var groupUsers []model.GroupUser
	err := db.Preload("User").Where("group_id = ?", groupID).Find(&groupUsers).Error
	return groupUsers, err
}
