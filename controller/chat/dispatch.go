package chat

import (
	"chat_server/message"
	"chat_server/service"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func dispatch(msg *message.Message) {
	t := msg.Type
	if 10 <= t && t < 20 {
		friendMessageHandle(msg)
	} else if 20 <= t && t < 30 {
		groupMessageHandle(msg)
	} else if 50 <= t && t < 60 {
		friendNotifyHandle(msg)
	} else if 60 <= t && t < 70 {
		groupNotifyHandle(msg)
	} else {
		log.Error("消息类型错误")
	}
}

func readHandle(c *client) {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
		// 下线处理
		go userOfflineHandle(c.name)
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
		err = service.SaveMessage(msg)
		if err != nil {
			log.Error(err)
		}
		dispatch(msg)
	}
}

func writeHandle(c *client) {
	defer c.conn.Close()

	for msg := range c.send {
		data, err := proto.Marshal(msg)
		if err != nil {
			log.Error(err)
			continue
		}
		c.conn.WriteMessage(websocket.BinaryMessage, data)
	}
	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}
