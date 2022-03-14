package service

import "chat_server/model"

func GetUserById(id string) *model.User {
	return model.GetUserById(id)
}

func InsertUser(user *model.User) error {
	return model.InsertUser(user)
}
