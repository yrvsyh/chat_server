package chat

import (
	"chat_server/message"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type (
	client struct {
		name    string
		friends map[string]bool
		groups  map[uint]bool
		conn    *websocket.Conn
		send    chan *message.Message
		sendRaw chan []byte
	}

	friendsInfo struct {
		name    string
		friends map[string]bool
	}

	groupsInfo struct {
		name   string
		groups map[uint]bool
	}

	clientManager struct {
		clients       map[string]*client
		register      chan *client
		unregister    chan *client
		friendsChange chan *friendsInfo
		groupsChange  chan *groupsInfo
		// broadcast     chan []byte
	}
)

var manager = clientManager{
	clients:       make(map[string]*client),
	register:      make(chan *client, 16),
	unregister:    make(chan *client, 16),
	friendsChange: make(chan *friendsInfo, 16),
	groupsChange:  make(chan *groupsInfo, 16),
}

func init() {
	go manager.run()
}

func (m *clientManager) run() {
	for {
		select {
		case client := <-m.register:
			m.clients[client.name] = client
			log.Info("Register: ", client.name)
		case client := <-m.unregister:
			delete(m.clients, client.name)
			log.Info("Unregister: ", client.name)
		case friendsInfo := <-m.friendsChange:
			client, ok := m.clients[friendsInfo.name]
			if ok {
				client.friends = friendsInfo.friends
			}
			log.Info("FriendsChange: ", client.name)
		case groupsInfo := <-m.groupsChange:
			client, ok := m.clients[groupsInfo.name]
			if ok {
				client.groups = groupsInfo.groups
			}
			log.Info("GroupsChange: ", client.name)
		}
	}
}

func unregisterClient(c *client) {
	c.conn.Close()
	manager.unregister <- c
}

func RegisterClient(name string, conn *websocket.Conn) {
	client := &client{name: name, conn: conn, send: make(chan *message.Message, 128), sendRaw: make(chan []byte, 128)}
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
	client, ok := manager.clients[msg.To]
	if ok {
		client.send <- msg
	}
}

func UpdateUserFriendsInfo(name string, friends map[string]bool) {
	friendsInfo := &friendsInfo{name: name, friends: friends}
	manager.friendsChange <- friendsInfo
}

func UpdateUserGroupsInfo(name string, groups map[uint]bool) {
	groupsInfo := &groupsInfo{name: name, groups: groups}
	manager.groupsChange <- groupsInfo
}
