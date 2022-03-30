package chat_service

import (
	"chat_server/message"

	log "github.com/sirupsen/logrus"
)

func userOnlineHandle(c *client) bool {
	// 初始化manager内的friends和groups信息
	friendsSet, err := userService.GetUserFriendIDSet(c.id)
	if err != nil {
		log.Error(err)
		return false
	}
	c.friends = friendsSet

	groupsSet, err := userService.GetUserGroupIDSet(c.id)
	if err != nil {
		log.Error(err)
		return false
	}
	c.groups = groupsSet

	// 对在线好友广播上线通知
	for friend := range c.friends {
		c, ok := manager.clients[friend]
		// 好友在线
		if ok {
			msg := &message.Message{
				Type: message.Type_FRIEND_ONLINE,
				From: c.id,
				To:   friend,
			}
			c.send <- msg
		}
	}
	return true
}

func userOfflineHandle(c *client) {
	// 对在线好友广播下线通知
	for friend := range c.friends {
		c, ok := manager.clients[friend]
		// 好友在线
		if ok {
			msg := &message.Message{
				Type: message.Type_FRIEND_OFFLINE,
				From: c.id,
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
		sendMessage(msg.To, msg)
	}
}

func groupMessageHandle(msg *message.Message) {
	groupID := msg.To
	// 不能发给未加入的组
	_, ok := manager.clients[msg.From].groups[groupID]
	if ok {
		groupMembers := getGroupMembers(groupID)
		for id := range groupMembers {
			sendMessage(id, msg)
		}
	}
}

// func sendToGroup(group string, msg *message.Message) {
// 	key := config.RedisChannelGroupMessageKeyPrefix + group
// 	data, err := proto.Marshal(msg)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// 	err = RDB.Publish(key, data).Err()
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// }

// func recvFromGroup(id uint32) {
// 	groups := manager.clients[id].groups
// 	send := manager.clients[id].sendRaw
// 	for group := range groups {
// 		groupId := strconv.FormatUint(uint64(group), 10)
// 		go func(group string, send chan []byte) {
// 			ch := RDB.Subscribe(config.RedisChannelGroupMessageKeyPrefix + group).Channel()
// 			for msg := range ch {
// 				send <- []byte(msg.Payload)
// 			}
// 		}(groupId, send)
// 	}
// }
