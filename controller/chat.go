package controller

import (
	"chat_server/database"
	"chat_server/message"
	"chat_server/middleware"
	"chat_server/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type (
	clientManager struct {
		clients    map[string]*client
		register   chan *client
		unregister chan *client
		broadcast  chan []byte
	}

	client struct {
		id   string
		conn *websocket.Conn
		send chan []byte
	}
)

var (
	upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	manager = clientManager{
		clients:    make(map[string]*client),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
)

func init() {
	go manager.run()
}

func (m *clientManager) run() {
	for {
		select {
		case client := <-m.register:
			m.clients[client.id] = client
		case client := <-m.unregister:
			delete(m.clients, client.id)
		}
	}
}

func ChatHandle(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	tokenString := middleware.GetToken(c)
	claims, _ := middleware.ParseToken(tokenString)
	id := claims.Name

	client := &client{id: id, conn: ws, send: make(chan []byte)}
	manager.register <- client

	go readHandle(client)
	go writeHandle(client)

}

func readHandle(c *client) {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
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
			continue
		}
		log.Info(msg)

		id := msg.To

		client, ok := manager.clients[id]
		if !ok {
			// 目标用户不存在
			if !service.CheckUserExist(id) {
				continue
			}
			// 缓存离线信息
			database.RDB.LPush(id, msg)
		}

		// 重新构造消息内容
		msg.Id = time.Now().Unix()
		msg.To = msg.From
		msg.From = c.id

		data, err = proto.Marshal(msg)
		if err != nil {
			log.Error(err)
			continue
		}
		// 转发消息
		client.send <- data
	}
}

func writeHandle(c *client) {
	defer c.conn.Close()

	for data := range c.send {
		c.conn.WriteMessage(websocket.BinaryMessage, data)
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
