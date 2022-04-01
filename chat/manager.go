package chat

import (
	"chat_server/message"
	"sync"
	"time"

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
	msg.Id = time.Now().UnixMicro()
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
	groupMemberSetAny, _ := m.groupMembersMap.LoadOrStore(groupID, mapset.NewSet[uint32]())
	return groupMemberSetAny.(mapset.Set[uint32])
}

// 发送到指定用户
func (m *ChatManager) sendMessage(to uint32, msg *message.Message) {
	if client := m.getClient(to); client != nil {
		client.sendMessage(msg)
	}
}
