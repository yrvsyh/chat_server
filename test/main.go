package main

import (
	"chat_server/message"
	"chat_server/utils"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/shiena/ansicolor"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

func init() {
	log.StandardLogger().SetFormatter(&log.TextFormatter{ForceColors: true})
	log.StandardLogger().SetOutput(ansicolor.NewAnsiColorWriter(os.Stdout))
}

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
	header := http.Header{}
	header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDczMTg2NzQsImp0aSI6IjNlZTAyZDRlLWQxOTYtNDA4MS1hNTczLTg2NzYwMmJmNjRiZiIsImlhdCI6MTY0NzMxNTA3NCwibmFtZSI6Inl6eSJ9.G9FHL1Zn1eW9GULC2lXq4RNmRwJHUzQVMJgh-xgFk_A")
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws/chat", header)
	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	msg := &message.Message{}
	msg.Id = time.Now().Unix()
	msg.Type = 0
	msg.To = "yzy"
	msg.From = "yzy"
	msg.Content = make([]byte, 16)
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Error(err)
		return
	}

	err = conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		log.Error(err)
		return
	}
	_, data, err = conn.ReadMessage()
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			log.Error(err)
		}
		return
	}

	msg = &message.Message{}
	err = proto.Unmarshal(data, msg)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info(msg)

	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func main() {
	// testRedis()
	// testProto()
	testChat()
}
