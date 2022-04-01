package service

import "chat_server/database"

var (
	db = database.DB

	userService    = UserService{}
	groupService   = GroupService{}
	messageService = MessageService{}
)
