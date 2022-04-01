package chat

import (
	"chat_server/message"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
)

type ChatManager struct {
	clientsMap      *sync.Map // map[uint32]*client
	groupMembersMap *sync.Map // map[uint32]mapset.Set
}

func NewChatManager() *ChatManager {
	return &ChatManager{
		clientsMap:      &sync.Map{},
		groupMembersMap: &sync.Map{},
	}
}

func (m *ChatManager) Register(id uint32, conn *websocket.Conn) {
	client := NewClient(m, id, conn)
	m.clientsMap.Store(id, client)
	client.Init()
}

// 服务端主动发送通知
func (m *ChatManager) SendMessage(msg *message.Message) {
	switch msg.Type {
	case message.Type_FRIEND_REQUEST, message.Type_FRIEND_ACCEPT, message.Type_FRIEND_DISBAND:
		messageService.SaveUserMessage(msg)
	case message.Type_GROUP_REQUEST, message.Type_GROUP_ACCEPT, message.Type_GROUP_DISBAND, message.Type_GROUP_USER_CHANGE:
		messageService.SaveGroupMessage(msg)
	}
	m.sendMessage(msg.To, msg)
}

func (m *ChatManager) UpdateUserFriendsInfo(userID uint32, friendID uint32, isFriend bool) {
	clientAny, ok := m.clientsMap.Load(userID)
	if ok {
		client := clientAny.(*client)
		if isFriend {
			client.friendSet.Add(friendID)
		} else {
			client.friendSet.Remove(friendID)
		}
	}
}

func (m *ChatManager) UpdateUserGroupsInfo(userID uint32, groupID uint32, inGroup bool) {
	clientAny, ok := m.clientsMap.Load(userID)
	if ok {
		client := clientAny.(*client)
		groupMembers := m.getGroupMembers(groupID)
		if inGroup {
			client.groupSet.Add(groupID)
			groupMembers.Add(userID)
		} else {
			client.groupSet.Remove(groupID)
			groupMembers.Remove(userID)
		}
	}
}

func (m *ChatManager) getClient(id uint32) *client {
	var ret *client
	clientAny, ok := m.clientsMap.Load(id)
	if ok {
		c, ok := clientAny.(*client)
		if ok {
			ret = c
		}
	}
	return ret
}

func (m *ChatManager) getGroupMembers(groupID uint32) mapset.Set[uint32] {
	groupMemberSetAny, ok := m.groupMembersMap.Load(groupID)
	if !ok {
		m.groupMembersMap.Store(groupID, mapset.NewSet[uint32]())
		groupMemberSetAny, _ = m.groupMembersMap.Load(groupID)
	}
	return groupMemberSetAny.(mapset.Set[uint32])
}

// 转发消息
func (m *ChatManager) sendMessage(id uint32, msg *message.Message) {
	if client := m.getClient(id); client != nil {
		client.send <- msg
	}
}
