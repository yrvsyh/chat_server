package main

import (
	"chat_server/database"
	"chat_server/model"
	"fmt"
)

func main() {
	err := database.DB.AutoMigrate(model.User{}, model.UserFriend{}, model.UserMessage{}, model.Group{}, model.GroupUser{}, model.GroupMessage{})
	fmt.Println(err)
}
