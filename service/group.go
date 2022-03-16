package service

import (
	"chat_server/model"
)

func GetGroupMemberNameList(id uint) ([]string, error) {
	var ret []string
	var userGroups []model.UserGroups
	err := db.Model(&model.UserGroups{}).Where("group_id = ?", id).Find(&userGroups).Error
	for _, userGroup := range userGroups {
		ret = append(ret, userGroup.UserName)
	}
	return ret, err
}

func GetGroupMembers(id uint) ([]*model.User, error) {
	group := &model.Group{}
	err := db.Preload("Members").First(group, id).Error
	return group.Members, err
}
