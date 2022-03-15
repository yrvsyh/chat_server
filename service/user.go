package service

import "chat_server/model"

func GetUserById(id string) *model.User {
	return model.GetUserById(id)
}

func InsertUser(user *model.User) error {
	return model.InsertUser(user)
}

func CheckUserExist(id string) bool {
	return model.GetUserById(id) != nil
}

func GetFriendsList(id string) []model.User {
	return make([]model.User, 0)
}
