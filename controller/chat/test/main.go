package main

import (
	"chat_server/message"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	fmt.Println("name password")
	var name string
	var password string
	_, _ = fmt.Scanf("%s %s\n", &name, &password)

	body := fmt.Sprintf("username=%s&password=%s", name, password)
	resp, err := http.Post("http://127.0.0.1:8080/auth/login", "application/x-www-form-urlencoded", strings.NewReader(body))
	if err != nil {
		log.Error(err)
		return
	}

	data, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(data))

	req := &http.Request{Header: make(http.Header)}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth" {
			req.AddCookie(cookie)
		}
	}

	//fmt.Printf("%+v\n", req.Header)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws/chat", req.Header)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	go func(conn *websocket.Conn) {
		for {
			_, data, err = conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Error(err)
				}
				return
			}

			msg := &message.Message{}
			err = proto.Unmarshal(data, msg)
			if err != nil {
				log.Error(err)
				return
			}
			fmt.Println(msg.Type, string(msg.Content))
		}
	}(conn)

	for {
		fmt.Printf(">>>")
		var toName string
		var msgContent string
		_, _ = fmt.Scanf("%s %s\n", &toName, &msgContent)
		if toName == "" {
			continue
		}

		msg := &message.Message{}
		msg.Id = time.Now().UnixNano()
		msg.Type = int32(message.MessageType_FRIEND_TEXT)
		msg.From = name
		msg.To = toName
		msg.Content = []byte(msgContent)
		data, err = proto.Marshal(msg)
		if err != nil {
			log.Error(err)
			break
		}

		err = conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Error(err)
			break
		}
	}

	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
