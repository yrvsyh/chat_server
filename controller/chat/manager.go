package chat

import (
	"chat_server/message"
	"github.com/gorilla/websocket"
)

type (
	client struct {
		name    string
		friends map[string]bool
		groups  map[uint]bool
		conn    *websocket.Conn
		send    chan *message.Message
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
		broadcast     chan []byte
		friendsChange chan *friendsInfo
		groupsChange  chan *groupsInfo
	}
)

var manager = clientManager{
	clients:       make(map[string]*client),
	register:      make(chan *client),
	unregister:    make(chan *client),
	friendsChange: make(chan *friendsInfo),
	groupsChange:  make(chan *groupsInfo),
}

func init() {
	go manager.run()
}

func (m *clientManager) run() {
	for {
		select {
		case client := <-m.register:
			m.clients[client.name] = client
		case client := <-m.unregister:
			delete(m.clients, client.name)
		case friendsInfo := <-m.friendsChange:
			client, ok := m.clients[friendsInfo.name]
			if ok {
				client.friends = friendsInfo.friends
			}
		case groupsInfo := <-m.groupsChange:
			client, ok := m.clients[groupsInfo.name]
			if ok {
				client.groups = groupsInfo.groups
			}
		}
	}
}

func RegisterClient(name string, conn *websocket.Conn) {
	client := &client{name: name, conn: conn, send: make(chan *message.Message)}
	manager.register <- client

	go readHandle(client)
	go writeHandle(client)

	// 上线处理
	go userOnlineHandle(name)
}

func UpdateUserFriendsInfo(name string, friends map[string]bool) {
	friendsInfo := &friendsInfo{name: name, friends: friends}
	manager.friendsChange <- friendsInfo
}

func UpdateUserGroupsInfo(name string, groups map[uint]bool) {
	groupsInfo := &groupsInfo{name: name, groups: groups}
	manager.groupsChange <- groupsInfo
}
