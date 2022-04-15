package service

import "chat_server/database"

var (
	db = database.GetMysqlInstance()

	userService    = UserService{}
	groupService   = GroupService{}
	messageService = MessageService{}
)
