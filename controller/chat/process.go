package chat

import (
	"chat_server/message"
	"chat_server/service"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func userOnlineHandle(name string) {
	// 对在线好友广播上线通知
	friends, err := service.GetUserFriendNameList(name)
	if err != nil {
		log.Error(err)
	} else {
		for _, friend := range friends {
			c, ok := manager.clients[friend]
			// 好友在线
			if ok {
				msg := &message.Message{
					Type: int32(message.MessageType_FRIEND_ONLINE),
					From: name,
					To:   friend,
				}
				c.send <- msg
			}
		}
	}
}

func userOfflineHandle(name string) {
	// 对在线好友广播下线通知
	friends, err := service.GetUserFriendNameList(name)
	if err != nil {
		log.Error(err)
	} else {
		for _, friend := range friends {
			c, ok := manager.clients[friend]
			// 好友在线
			if ok {
				msg := &message.Message{
					Type: int32(message.MessageType_FRIEND_OFFLINE),
					From: name,
					To:   friend,
				}
				c.send <- msg
			}
		}
	}
}

func friendMessageHandle(msg *message.Message) {
	name := msg.To
	c, ok := manager.clients[name]
	if ok {
		c.send <- msg
	}
}

func groupMessageHandle(msg *message.Message) {
	id, err := strconv.ParseUint(msg.To, 0, 32)
	if err != nil {
		log.Error(err)
		return
	}
	membersName, err := service.GetGroupMemberNameList(uint(id))
	for _, name := range membersName {
		c, ok := manager.clients[name]
		if ok {
			c.send <- msg
		}
	}
}

func friendNotifyHandle(msg *message.Message) {}

func groupNotifyHandle(msg *message.Message) {}
