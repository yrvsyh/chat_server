package chat

import (
	"chat_server/message"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func dispatch(msg *message.Message) {
	t := msg.Type
	if 10 < t && t < 20 {
		friendMessageHandle(msg)
	} else if 20 <= t && t < 30 {
		groupMessageHandle(msg)
	} else {
		log.Error("消息类型错误")
	}
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
		log.Info(msg)
		// 记录至数据库
		err = messageService.SaveMessage(msg)
		if err != nil {
			log.Error(err)
		}
		dispatch(msg)
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
