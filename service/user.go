package service

import "chat_server/model"

func RegisterUser(auth *model.UserAuth) error {
	user := &model.User{
		Username: auth.Username,
		UserAuth: auth,
	}
	return db.Create(user).Error
}

func GetUserAuthByName(username string) (*model.User, error) {
	user := &model.User{}
	err := db.Preload("UserAuth").Where("username = ?", username).First(user).Error
	return user, err
}

func GetUserById(id uint) (*model.User, error) {
	user := &model.User{}
	err := db.First(user, id).Error
	return user, err
}

func GetUserByName(username string) (*model.User, error) {
	user := &model.User{}
	err := db.Where("username = ?", username).First(user).Error
	return user, err
}

func GetUserFriendsById(id uint) ([]*model.UserFriends, error) {
	user := &model.User{}
	err := db.Preload("Friends").First(user, id).Error
	return user.Friends, err
}

func GetUserFriendsDetailById(id uint) ([]*model.UserFriends, error) {
	user := &model.User{}
	err := db.Preload("Friends.Friend").First(user, id).Error
	return user.Friends, err
}

func CheckUserExistByName(username string) bool {
	_, err := GetUserByName(username)
	return err == nil
}
