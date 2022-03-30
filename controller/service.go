package controller

import "chat_server/service"

var (
	userService    = service.UserService{}
	groupService   = service.GroupService{}
	messageService = service.MessageService{}
)
