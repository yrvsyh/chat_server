package main

import (
	"chat_server/message"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
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
	fmt.Println(string(data))
	js := struct {
		Code int                    `json:"code,omitempty"`
		Msg  string                 `json:"msg,omitempty"`
		Data map[string]interface{} `json:"data,omitempty"`
	}{}
	json.Unmarshal(data, &js)
	fmt.Println(js)
	id := uint32(js.Data["id"].(float64))

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
		var to uint32
		var msgContent string
		var t int
		_, _ = fmt.Scanf("%d %d %s\n", &t, &to, &msgContent)

		msg := &message.Message{}
		msg.Id = time.Now().UnixNano()
		if t == 0 {
			msg.Type = message.Type_FRIEND_TEXT
		} else {
			msg.Type = message.Type_GROUP_TEXT
		}
		msg.From = id
		msg.To = to
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
