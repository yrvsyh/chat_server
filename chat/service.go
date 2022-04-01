package chat

import (
	"chat_server/database"
	"chat_server/service"
)

var (
	userService    = service.UserService{}
	groupService   = service.GroupService{}
	messageService = service.MessageService{}

	rdb = database.RDB
)
