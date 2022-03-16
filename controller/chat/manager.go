package chat

import (
	"chat_server/message"
	"github.com/gorilla/websocket"
)

type (
	client struct {
		name string
		conn *websocket.Conn
		send chan *message.Message
	}

	clientManager struct {
		clients    map[string]*client
		register   chan *client
		unregister chan *client
		broadcast  chan []byte
	}
)

var manager = clientManager{
	clients:    make(map[string]*client),
	register:   make(chan *client),
	unregister: make(chan *client),
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
