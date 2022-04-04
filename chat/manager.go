package chat

import (
	"chat_server/message"
	"chat_server/utils"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
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
	// 单客户端登录
	_, ok := m.clientsMap.Load(id)
	if ok {
		log.Info("Signed in on another client")
		conn.WriteMessage(websocket.CloseMessage, []byte("Signed in on another client"))
		conn.Close()
		return
	}

	client := NewClient(m, id, conn)
	m.clientsMap.Store(id, client)
	client.Init()
}

// 服务端主动发送通知
func (m *ChatManager) SendMessage(msg *message.Message) {
	msg.Id = utils.GenMsgID()
	switch msg.Type {
	case message.Type_FRIEND_REQUEST, message.Type_FRIEND_ACCEPT, message.Type_FRIEND_DISBAND:
		messageService.SaveMessage(msg)
	case message.Type_GROUP_REQUEST, message.Type_GROUP_ACCEPT, message.Type_GROUP_DISBAND, message.Type_GROUP_USER_CHANGE:
		messageService.SaveMessage(msg)
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
	groupMemberSetAny, ok := m.groupMembersMap.LoadOrStore(groupID, mapset.NewSet[uint32]())
	groupMemberSet := groupMemberSetAny.(mapset.Set[uint32])
	if !ok {
		members, err := groupService.GetGroupMemberSet(groupID)
		if err != nil {
			log.Error(err)
		}
		// TODO 使用mapset.Set
		for userID := range members {
			groupMemberSet.Add(userID)
		}
	}
	return groupMemberSet
}

// 发送到指定用户
func (m *ChatManager) sendMessage(to uint32, msg *message.Message) {
	if client := m.getClient(to); client != nil {
		client.sendMessage(msg)
	}
}
