package main

import (
	"chat_server/database"
	"chat_server/model"
	"fmt"
)

func main() {
	db := database.GetMysqlInstance()
	if err := db.AutoMigrate(model.User{}, model.UserFriend{}, model.UserMessage{}, model.Group{}, model.GroupUser{}, model.GroupMessage{}); err != nil {
		fmt.Println(err)
	}
}
