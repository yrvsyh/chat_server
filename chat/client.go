package chat

import (
	"chat_server/message"
	"chat_server/utils"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type pendingMessage struct {
	msg   *message.Message
	retry time.Duration
	timer *time.Timer
}

type client struct {
	m         *ChatManager
	id        uint32
	conn      *websocket.Conn
	friendSet mapset.Set[uint32]
	groupSet  mapset.Set[uint32]
	send      chan *message.Message
	sendRaw   chan []byte
	close     chan struct{}

	pendingMap *sync.Map // map[uint64]*pendingMessage
}

func NewClient(m *ChatManager, id uint32, conn *websocket.Conn) *client {
	return &client{
		m:         m,
		id:        id,
		conn:      conn,
		friendSet: mapset.NewSet[uint32](),
		groupSet:  mapset.NewSet[uint32](),
		send:      make(chan *message.Message, 16),
		sendRaw:   make(chan []byte, 16),
		close:     make(chan struct{}),

		pendingMap: &sync.Map{},
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
	c.pendingMap.Range(func(key, value any) bool {
		pm := value.(*pendingMessage)
		pm.timer.Stop()
		c.pendingMap.Delete(key)
		return true
	})

	c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	c.conn.Close()

	c.m.clientsMap.Delete(c.id)
}

func (c *client) readHandle() {
	defer func() {
		// 下线处理
		c.close <- struct{}{}
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Error(err)
			}
			break
		}

		// 解析protobuf
		msg := &message.Message{}
		err = proto.Unmarshal(data, msg)
		if err != nil {
			logrus.Error(err)
			continue
		}
		// 不是这个客户端的消息
		if msg.From != c.id {
			continue
		}

		switch msg.Type {
		case message.Type_Acknowledge:
			c.AckHandle(msg.Id)
		case message.Type_FRIEND_TEXT, message.Type_FRIEND_IMAGE, message.Type_FRIEND_FILE:
			msg.Id = utils.GenMsgID()
			messageService.SaveMessage(msg)
			c.friendMessageHandle(msg)
		case message.Type_GROUP_TEXT, message.Type_GROUP_IMAGE, message.Type_GROUP_FILE:
			msg.Id = utils.GenMsgID()
			messageService.SaveMessage(msg)
			c.groupMessageHandle(msg)
		default:
			logrus.Error("message type error")
		}
	}
}

func (c *client) writeHandle() {
	defer func() {
		// 下线处理
		c.userOfflineHandle()
		c.unregister()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			data, err := proto.Marshal(msg)
			if err != nil {
				logrus.Error(err)
				continue
			}
			if c.conn.WriteMessage(websocket.BinaryMessage, data) != nil {
				return
			}
		case data, ok := <-c.sendRaw:
			if !ok {
				return
			}
			if c.conn.WriteMessage(websocket.BinaryMessage, data) != nil {
				return
			}
		case <-c.close:
			return
		}
	}
}

func (c *client) AckHandle(msgID int64) {
	value, ok := c.pendingMap.LoadAndDelete(msgID)
	if ok {
		pm := value.(*pendingMessage)

		// 关闭Ack超时计时器
		pm.timer.Stop()

		msg := pm.msg
		if message.CheckMessageType(msg) == message.MESSAGE_TYPE_NOTIFY {
			return
		}
		// 客户端已收到消息
		messageService.UpdateMessageState(msg, message.State_CLIENT_RECV)
		// 更新User最后收到的msgID
		// 更新这个客户端的信息
		msg.From = c.id
		if err := messageService.UpdateLastMsgID(msg); err != nil {
			logrus.Error(err)
		}
	}
}

func (c *client) ReplyAck(msg *message.Message) {
	ack := &message.Message{
		Id:      msg.Id,
		LocalId: msg.LocalId,
		Type:    message.Type_Acknowledge,
		From:    0,
		To:      c.id,
	}
	c.send <- ack
}

// 发送消息到当前客户端
func (c *client) sendMessage(msg *message.Message) {
	// 加入已发送队列
	pm := &pendingMessage{msg: msg, retry: time.Second}
	c.pendingMap.Store(msg.Id, pm)
	pm.timer = time.AfterFunc(time.Second, func() {
		// 未收到Ack, 重发消息
		time.Sleep(pm.retry)
		pm.timer.Reset(time.Second)
		c.send <- msg
		pm.retry *= 2
	})

	// 等待Ack
	messageService.UpdateMessageState(msg, message.State_WAIT_ACK)

	c.send <- msg
}
