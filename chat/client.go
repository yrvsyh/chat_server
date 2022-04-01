package chat

import (
	"chat_server/message"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type client struct {
	m         *ChatManager
	id        uint32
	conn      *websocket.Conn
	friendSet mapset.Set[uint32]
	groupSet  mapset.Set[uint32]
	send      chan *message.Message
	sendRaw   chan []byte
}

func NewClient(m *ChatManager, id uint32, conn *websocket.Conn) *client {
	return &client{
		m:         m,
		id:        id,
		conn:      conn,
		friendSet: mapset.NewSet[uint32](),
		groupSet:  mapset.NewSet[uint32](),
	}
}

func (c *client) Init() {
	go c.writeHandle()
	// 上线处理
	if !c.userOnlineHandle() {
		// 上线初始化出错后直接断开链接, 客户端自动重连
		c.unregister()
		return
	}

	go c.readHandle()
}

func (c *client) unregister() {
	id := c.id
	c.m.clientsMap.Delete(id)
	c.conn.Close()
}



func (c *client) readHandle() {
	defer func() {
		// 下线处理
		c.userOfflineHandle()
		c.unregister()
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
		// 不是这个客户端的消息
		if msg.From != c.id {
			continue
		}

		// TODO 回复Ack

		switch msg.Type {
		case message.Type_Acknowledge:
		case message.Type_FRIEND_TEXT, message.Type_FRIEND_IMAGE, message.Type_FRIEND_FILE:
			messageService.SaveUserMessage(msg)
			c.friendMessageHandle(msg)
		case message.Type_GROUP_TEXT, message.Type_GROUP_IMAGE, message.Type_GROUP_FILE:
			messageService.SaveGroupMessage(msg)
			c.groupMessageHandle(msg)
		}
	}
}

func (c *client) writeHandle() {
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
