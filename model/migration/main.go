package main

import (
	"chat_server/database"
	"chat_server/model"
	"fmt"
)

func main() {
	database.InitMysql()
	if err := database.DB.AutoMigrate(model.User{}, model.UserFriend{}, model.UserMessage{}, model.Group{}, model.GroupUser{}, model.GroupMessage{}); err != nil {
		fmt.Println(err)
	}
}
