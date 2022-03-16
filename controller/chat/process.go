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
	friends := manager.clients[name].friends
	for friend, _ := range friends {
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

func friendMessageHandle(msg *message.Message) {
	// 不能发给陌生人
	friends := manager.clients[msg.From].friends
	_, ok := friends[msg.To]
	if ok {
		c, ok := manager.clients[msg.To]
		// 直接转发给在线用户
		if ok {
			c.send <- msg
		}
	}
}

func groupMessageHandle(msg *message.Message) {
	id64, err := strconv.ParseUint(msg.To, 0, 32)
	id := uint(id64)
	if err != nil {
		log.Error(err)
		return
	}
	// 不能发给未加入的组
	groups := manager.clients[msg.From].groups
	_, ok := groups[id]
	if ok {
		membersName, err := service.GetGroupMemberNameList(id)
		if err != nil {
			log.Error(err)
			return
		}
		for _, name := range membersName {
			c, ok := manager.clients[name]
			// 直接转发给在线用户
			if ok {
				c.send <- msg
			}
		}
	}
}

func friendNotifyHandle(msg *message.Message) {}

func groupNotifyHandle(msg *message.Message) {}
