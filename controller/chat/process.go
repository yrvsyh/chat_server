package chat

import (
	"chat_server/config"
	"chat_server/message"
	"strconv"

	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func userOnlineHandle(c *client) bool {
	// 初始化manager内的friends和groups信息
	friendsSet, err := userService.GetUserFriendNameSet(c.name)
	if err != nil {
		log.Error(err)
		return false
	}
	UpdateUserFriendsInfo(c.name, friendsSet)

	groupsSet, err := userService.GetUserGroupNameSet(c.name)
	if err != nil {
		log.Error(err)
		return false
	}
	UpdateUserGroupsInfo(c.name, groupsSet)

	// 对在线好友广播上线通知
	friends := manager.clients[c.name].friends
	for friend := range friends {
		c, ok := manager.clients[friend]
		// 好友在线
		if ok {
			msg := &message.Message{
				Type: message.MessageType_FRIEND_ONLINE,
				From: c.name,
				To:   friend,
			}
			c.send <- msg
		}
	}
	return true
}

func userOfflineHandle(c *client) {
	// 对在线好友广播下线通知
	friends := c.friends
	for friend := range friends {
		c, ok := manager.clients[friend]
		// 好友在线
		if ok {
			msg := &message.Message{
				Type: message.MessageType_FRIEND_OFFLINE,
				From: c.name,
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
	id64, err := strconv.ParseUint(msg.To, 10, 32)
	id := uint(id64)
	if err != nil {
		log.Error(err)
		return
	}
	// 不能发给未加入的组
	groups := manager.clients[msg.From].groups
	_, ok := groups[id]
	if ok {
		membersName, err := groupService.GetGroupMemberNameList(id)
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

func sendToGroup(group string, msg *message.Message) {
	key := config.RedisChannelGroupMessageKeyPrefix + group
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}
	err = RDB.Publish(key, data).Err()
	if err != nil {
		log.Error(err)
		return
	}
}

func recvFromGroup(name string) {
	groups := manager.clients[name].groups
	send := manager.clients[name].sendRaw
	for group := range groups {
		groupId := strconv.FormatUint(uint64(group), 10)
		go func(group string, send chan []byte) {
			ch := RDB.Subscribe(config.RedisChannelGroupMessageKeyPrefix + group).Channel()
			for msg := range ch {
				send <- []byte(msg.Payload)
			}
		}(groupId, send)
	}
}

func testSet() {
	RDB.SMembers("key").Result()
}
