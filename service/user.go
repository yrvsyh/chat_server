package service

import (
	"chat_server/model"

	"gorm.io/gorm"
)

type UserService struct{}

func (UserService) RegisterUser(auth *model.UserAuth) error {
	user := &model.User{
		Name:     auth.UserName,
		UserAuth: auth,
	}
	return db.Create(user).Error
}

func (UserService) GetUserAuthByName(name string) (*model.User, error) {
	user := &model.User{}
	err := db.Preload("UserAuth").First(user, "name = ?", name).Error
	return user, err
}

func (UserService) GetUserByName(name string) (*model.User, error) {
	user := &model.User{}
	err := db.First(user, "name = ?", name).Error
	return user, err
}

func (UserService) UpdateUser(user *model.User) error {
	return db.Model(user).Updates(user).Error
}

func (UserService) UpdateUserAvatarByName(name string, path string) error {
	return db.Model(&model.User{}).Where("name = ?", name).Update("avatar", path).Error
}

func (UserService) GetUserFriendNameList(name string) ([]string, error) {
	var ret []string
	var userFriends []model.UserFriends
	err := db.Model(&model.UserFriends{}).Where("user_name = ?", name).Find(&userFriends).Error
	for _, userFriend := range userFriends {
		ret = append(ret, userFriend.FriendName)
	}
	return ret, err
}

func (UserService) GetUserFriendNameSet(name string) (map[string]bool, error) {
	ret := make(map[string]bool)
	var userFriends []model.UserFriends
	err := db.Model(&model.UserFriends{}).Where("user_name = ?", name).Find(&userFriends).Error
	for _, userFriend := range userFriends {
		ret[userFriend.FriendName] = true
	}
	return ret, err
}

func (UserService) GetUserGroupNameSet(name string) (map[uint]bool, error) {
	ret := make(map[uint]bool)
	var userGroups []model.UserGroups
	err := db.Model(&model.UserGroups{}).Where("user_name = ?", name).Find(&userGroups).Error
	for _, userGroup := range userGroups {
		ret[userGroup.GroupID] = true
	}
	return ret, err
}

func (UserService) GetUserFriendsByName(name string) ([]*model.UserFriends, error) {
	user := &model.User{}
	err := db.Preload("Friends", "accept = ?", true).First(user, "name = ?", name).Error
	return user.Friends, err
}

func (UserService) GetUserFriendsDetailByName(name string) ([]*model.UserFriends, error) {
	user := &model.User{}
	err := db.Preload("Friends", "accept = ?", true).Preload("Friends.Friend").First(user, "name = ?", name).Error
	return user.Friends, err
}

func (UserService) GetUserFriendDetailByFriendName(name string, friendName string) (*model.UserFriends, error) {
	userFriens := &model.UserFriends{}
	err := db.Model(userFriens).First(userFriens, "user_name = ? and friend_name = ?", name, friendName).Error
	return userFriens, err
}

func (UserService) UpdateUserFriend(userFriend *model.UserFriends) error {
	return db.Save(userFriend).Error
}

func (s UserService) CheckUserExistByName(name string) bool {
	_, err := s.GetUserByName(name)
	return err == nil
}

func (UserService) AddUserFriend(name string, friendName string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		user := &model.User{}
		err := db.First(user, "name = ?", name).Error
		if err != nil {
			return err
		}

		friend := &model.User{}
		err = tx.First(friend, "name = ?", friendName).Error
		if err != nil {
			return err
		}

		// 双向好友, 默认为未同意
		err = db.Model(user).Association("Friends").Append(&model.UserFriends{FriendName: friend.Name})
		err = db.Model(friend).Association("Friends").Append(&model.UserFriends{FriendName: user.Name})
		if err != nil {
			return err
		}

		return nil
	})
}

func (UserService) AcceptUserFriend(name string, friendName string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		userFriens1 := &model.UserFriends{UserName: name, FriendName: friendName}
		userFriens2 := &model.UserFriends{UserName: friendName, FriendName: name}
		// err := tx.Model(userFriens).First(userFriens, "user_name = ? and friend_name = ?", name, friendName).Error
		// if err != nil {
		// 	return err
		// }
		// userFriens.Accept = true
		err := tx.Model(userFriens1).Update("accept", true).Error
		err = tx.Model(userFriens2).Update("accept", true).Error

		if err != nil {
			return err
		}

		return nil
	})
}
