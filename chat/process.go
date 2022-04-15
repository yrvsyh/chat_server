package chat

import (
	"chat_server/message"
	"chat_server/utils"

	log "github.com/sirupsen/logrus"
)

func (c *client) userOnlineHandle() bool {
	// 初始化manager内的friends和groups信息
	friendsSet, err := userService.GetUserFriendIDSet(c.id)
	if err != nil {
		log.Error(err)
		return false
	}
	// TODO 使用mapset.Set
	for friend := range friendsSet {
		c.friendSet.Add(friend)
	}

	groupsSet, err := userService.GetUserGroupIDSet(c.id)
	if err != nil {
		log.Error(err)
		return false
	}
	for group := range groupsSet {
		c.groupSet.Add(group)
	}

	// 对在线好友广播上线通知
	c.friendSet.Each(func(friend uint32) bool {
		_, ok := c.m.clientsMap.Load(friend)
		if ok {
			msg := &message.Message{
				Id:   utils.GenMsgID(),
				Type: message.Type_FRIEND_ONLINE,
				From: c.id,
				To:   friend,
			}
			c.m.sendMessage(friend, msg)
		}
		return false
	})

	// TODO 发送状态为等待Ack的消息

	// 发送所有未接收的通知消息 包括重要的通知消息
	// FIXME 根据State筛选未发送的消息 排除等待Ack中的消息
	friendMessages, _ := messageService.GetFriendOfflineMessages(c.id)
	for _, friendMessage := range friendMessages {
		msg := messageService.LoadMessage(&friendMessage.Message)
		c.sendMessage(msg)
	}

	// FIXME 群消息等待用户拉取
	groupMessages, _ := messageService.GetGroupOfflineMessages(c.id)
	for _, groupMessage := range groupMessages {
		msg := messageService.LoadMessage(&groupMessage.Message)
		c.sendMessage(msg)
	}

	return true
}

func (c *client) userOfflineHandle() {
	// 对在线好友广播下线通知
	c.friendSet.Each(func(friend uint32) bool {
		_, ok := c.m.clientsMap.Load(friend)
		if ok {
			msg := &message.Message{
				Id:   utils.GenMsgID(),
				Type: message.Type_FRIEND_OFFLINE,
				From: c.id,
				To:   friend,
			}
			c.m.sendMessage(friend, msg)
		}
		return false
	})
}

func (c *client) friendMessageHandle(msg *message.Message) {
	friends := c.friendSet
	if friends.Contains(msg.To) {
		c.m.sendMessage(msg.To, msg)
		c.ReplyAck(msg)
	}
}

func (c *client) groupMessageHandle(msg *message.Message) {
	groupID := msg.To
	if c.groupSet.Contains(groupID) {
		groupMembers := c.m.getGroupMembers(groupID)
		groupMembers.Each(func(userID uint32) bool {
			c.m.sendMessage(userID, msg)
			return false
		})
		c.ReplyAck(msg)
	}
}
