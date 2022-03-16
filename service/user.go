package service

import "chat_server/model"

func RegisterUser(auth *model.UserAuth) error {
	user := &model.User{
		Name:     auth.UserName,
		UserAuth: auth,
	}
	return db.Create(user).Error
}

func GetUserAuthByName(name string) (*model.User, error) {
	user := &model.User{}
	err := db.Preload("UserAuth").First(user, "name = ?", name).Error
	return user, err
}

func GetUserByName(name string) (*model.User, error) {
	user := &model.User{}
	err := db.First(user, "name = ?", name).Error
	return user, err
}

func UpdateUser(user *model.User) error {
	return db.Model(user).Updates(user).Error
}

func UpdateUserAvatarByName(name string, path string) error {
	return db.Model(&model.User{}).Where("name = ?", name).Update("avatar", path).Error
}

func GetUserFriendNameList(name string) ([]string, error) {
	var ret []string
	var userFriends []model.UserFriends
	err := db.Model(&model.UserFriends{}).Where("user_name = ?", name).Find(&userFriends).Error
	for _, userFriend := range userFriends {
		ret = append(ret, userFriend.FriendName)
	}
	return ret, err
}

func GetUserFriendNameSet(name string) (map[string]bool, error) {
	var ret map[string]bool
	var userFriends []model.UserFriends
	err := db.Model(&model.UserFriends{}).Where("user_name = ?", name).Find(&userFriends).Error
	for _, userFriend := range userFriends {
		ret[userFriend.FriendName] = true
	}
	return ret, err
}

func GetUserFriendsByName(name string) ([]*model.UserFriends, error) {
	user := &model.User{}
	err := db.Preload("Friends", "accept = ?", true).First(user, "name = ?", name).Error
	return user.Friends, err
}

func GetUserFriendsDetailByName(name string) ([]*model.UserFriends, error) {
	user := &model.User{}
	err := db.Preload("Friends", "accept = ?", true).Preload("Friends.Friend").First(user, "name = ?", name).Error
	return user.Friends, err
}

func CheckUserExistByName(Name string) bool {
	_, err := GetUserByName(Name)
	return err == nil
}
