package main

import (
	"chat_server/message"
	"chat_server/utils"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func testRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		// log.Panic().Err(err)
		log.Panic(err)
	}

}

func testProto() {
	msg := message.Message{
		Id: 1,
	}
	ret, _ := proto.Marshal(&msg)
	str := utils.Bytes2BinStr(ret)
	fmt.Println(ret)
	fmt.Println(str)
}

func testChat() {
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws/chat", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	err = conn.WriteMessage(websocket.BinaryMessage, []byte("hello world"))
	if err != nil {
		log.Fatal(err)
	}
	mt, data, err := conn.ReadMessage()
	if err != nil {
		log.Fatal(err.(*websocket.CloseError))
	}
	log.WithFields(log.Fields{
		"Type": mt,
		"Data": string(data),
	}).Info("RECV")
	err = conn.WriteMessage(websocket.CloseMessage, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// testRedis()
	// testProto()
	testChat()
}
