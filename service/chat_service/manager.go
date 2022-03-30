package chat_service

import (
	"chat_server/message"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type (
	client struct {
		id      uint32
		friends map[uint32]bool
		groups  map[uint32]bool
		conn    *websocket.Conn
		send    chan *message.Message
		sendRaw chan []byte
	}

	friendInfo struct {
		userID   uint32
		friendID uint32
		isFriend bool
	}

	groupInfo struct {
		userID  uint32
		groupID uint32
		inGroup bool
	}

	clientManager struct {
		clients       map[uint32]*client
		groupMembers  map[uint32]map[uint32]bool
		register      chan *client
		unregister    chan *client
		friendsChange chan *friendInfo
		groupsChange  chan *groupInfo
		// broadcast     chan []byte
	}
)

var manager = clientManager{
	clients:       make(map[uint32]*client),
	groupMembers:  make(map[uint32]map[uint32]bool),
	register:      make(chan *client, 16),
	unregister:    make(chan *client, 16),
	friendsChange: make(chan *friendInfo, 16),
	groupsChange:  make(chan *groupInfo, 16),
}

func init() {
	go saveMessageHandle()
	go run()
}

func run() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client.id] = client
			log.Info("Register: ", client.id)
		case client := <-manager.unregister:
			delete(manager.clients, client.id)
			log.Info("Unregister: ", client.id)
		case friendsInfo := <-manager.friendsChange:
			client, ok := manager.clients[friendsInfo.userID]
			if ok {
				if friendsInfo.isFriend {
					client.friends[friendsInfo.friendID] = true
				} else {
					delete(client.friends, friendsInfo.friendID)
				}
				log.Info("FriendsChange: ", friendsInfo)
			}
		case groupsInfo := <-manager.groupsChange:
			client, ok := manager.clients[groupsInfo.userID]
			if ok {
				if groupsInfo.inGroup {
					client.groups[groupsInfo.groupID] = true
				} else {
					delete(client.groups, groupsInfo.groupID)
				}
				log.Info("GroupsChange: ", groupsInfo)
			}

			groupMember, ok := manager.groupMembers[groupsInfo.groupID]
			if ok {
				if groupsInfo.inGroup {
					groupMember[groupsInfo.userID] = true
				} else {
					delete(groupMember, groupsInfo.userID)
				}
				log.Info("GroupMembersChange: ", groupsInfo)
			}
		}
	}
}

func saveMessageHandle() {
}

func saveMessage(msg *message.Message) {
	log.Info("SAVE MESSAGE", msg)
}

func readHandle(c *client) {
	defer func() {
		// 下线处理
		userOfflineHandle(c)
		unregisterClient(c)
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error(err)
			}
			break
		}

		// 解析protobuf
		msg := &message.Message{}
		err = proto.Unmarshal(data, msg)
		if err != nil {
			log.Error(err)
			return
		}
		// 不是这个客户端的消息
		if msg.From != c.id {
			continue
		}

		switch msg.Type {
		case message.Type_Acknowledge:
		case message.Type_FRIEND_TEXT, message.Type_FRIEND_IMAGE, message.Type_FRIEND_FILE:
			messageService.SaveUserMessage(msg)
			friendMessageHandle(msg)
		case message.Type_GROUP_TEXT, message.Type_GROUP_IMAGE, message.Type_GROUP_FILE:
			messageService.SaveGroupMessage(msg)
			groupMessageHandle(msg)
		}
	}
}

func writeHandle(c *client) {
	defer func() {
		c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			data, err := proto.Marshal(msg)
			if err != nil {
				log.Error(err)
				continue
			}
			c.conn.WriteMessage(websocket.BinaryMessage, data)
		case data, ok := <-c.sendRaw:
			if !ok {
				return
			}
			c.conn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}

func unregisterClient(c *client) {
	c.conn.Close()
	manager.unregister <- c
}

func sendMessage(id uint32, msg *message.Message) {
	c, ok := manager.clients[id]
	if ok {
		// TODO ack处理
		c.send <- msg
	}
}

func getGroupMembers(groupID uint32) map[uint32]bool {
	groupMembers, ok := manager.groupMembers[groupID]
	if ok {
		return groupMembers
	} else {
		groupMembers, err := groupService.GetGroupMemberSet(groupID)
		if err != nil {
			return make(map[uint32]bool)
		}
		manager.groupMembers[groupID] = groupMembers
		return manager.groupMembers[groupID]
	}
}

func RegisterClient(id uint32, conn *websocket.Conn) {
	client := &client{id: id, conn: conn, send: make(chan *message.Message, 128), sendRaw: make(chan []byte, 128)}
	manager.register <- client

	go writeHandle(client)
	// 上线处理
	if !userOnlineHandle(client) {
		// 上线初始化出错后直接断开链接, 客户端自动重连
		unregisterClient(client)
		return
	}

	go readHandle(client)
}

func SendMessage(msg *message.Message) {
	switch msg.Type {
	case message.Type_FRIEND_REQUEST, message.Type_FRIEND_ACCEPT, message.Type_FRIEND_DISBAND:
		messageService.SaveUserMessage(msg)
	case message.Type_GROUP_REQUEST, message.Type_GROUP_ACCEPT, message.Type_GROUP_DISBAND, message.Type_GROUP_USER_CHANGE:
		messageService.SaveGroupMessage(msg)
	}
	sendMessage(msg.To, msg)
}

func UpdateUserFriendsInfo(userID uint32, friendID uint32, isFriend bool) {
	_, ok := manager.clients[userID]
	if ok {
		friendInfo := &friendInfo{userID: userID, friendID: friendID, isFriend: isFriend}
		manager.friendsChange <- friendInfo
	}
}

func UpdateUserGroupsInfo(userID uint32, groupID uint32, inGroup bool) {
	_, ok := manager.clients[userID]
	if ok {
		groupInfo := &groupInfo{userID: userID, groupID: groupID, inGroup: inGroup}
		manager.groupsChange <- groupInfo
	}
}
