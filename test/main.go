package main

import (
	"chat_server/database"
	"chat_server/message"
	"chat_server/model"
	"chat_server/service"
	"chat_server/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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

func testLogin() {
	resp, err := http.Post("http://127.0.0.1:8080/auth/login", "application/x-www-form-urlencoded", strings.NewReader("username=yzy&password=yuan"))
	if err != nil {
		log.Error(err)
		return
	}
	header := resp.Header
	log.Info(header)

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/auth/logout", nil)

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth_info" {
			req.AddCookie(cookie)
		}
	}

	log.Info(req.Cookies())

	client := &http.Client{}
	resp, _ = client.Do(req)
	data, _ := io.ReadAll(resp.Body)
	log.Info(string(data))
}

func testChat() {
	resp, err := http.Post("http://127.0.0.1:8080/auth/login", "application/x-www-form-urlencoded", strings.NewReader("username=yzy&password=yuan"))
	if err != nil {
		log.Error(err)
		return
	}

	data, _ := io.ReadAll(resp.Body)
	log.Info(string(data))

	req := &http.Request{Header: make(http.Header)}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "auth_info" {
			req.AddCookie(cookie)
		}
	}

	log.Infof("%+v", req.Header)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080/ws/chat", req.Header)
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
	data, err = proto.Marshal(msg)
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

func testDatabase() {
	db := database.DB

	user1 := &model.User{Name: "y", UserAuth: &model.UserAuth{Password: "x"}}
	user2 := &model.User{Name: "yy", UserAuth: &model.UserAuth{Password: "x"}}
	user3 := &model.User{Name: "yyy", UserAuth: &model.UserAuth{Password: "x"}}
	db.Create(user1)
	db.Create(user2)
	db.Create(user3)

	db.Where(user1, "name").First(user1)
	db.Where(user2, "name").First(user2)
	db.Where(user3, "name").First(user3)

	yzy := &model.User{}
	db.Where("name = ?", "yzy").First(yzy)

	//group := &model.Group{Name: "group"}
	//db.Create(group)

	db.Model(yzy).Association("Friends").Append(&model.UserFriends{FriendName: user1.Name, Remark: "yzy1"})
	db.Model(yzy).Association("Friends").Append(&model.UserFriends{FriendName: user2.Name, Remark: "yzy2"})
	db.Model(yzy).Association("Friends").Append(&model.UserFriends{FriendName: user3.Name, Remark: "yzy3"})

	db.Preload("Friends.Friend").Find(yzy)
	for _, friend := range user1.Friends {
		log.Infof("%+v", friend.Friend)
	}

	//db.Model(user1).Association("Groups").Append(&model.UserGroups{GroupID: group.ID, Remark: "group_remark"})
	//
	//db.Preload("Users").Find(group)
	//for _, user := range group.Users {
	//	log.Infof("%+v", user)
	//}

	////db.Where("username=?", user.Username).Delete(user)
	//userFriends1 := &model.UserFriends{UserID: user1.ID, FriendID: user2.ID, Remark: "yzy2"}
	//db.Create(userFriends1)
	//userFriends2 := &model.UserFriends{UserID: user1.ID, FriendID: user3.ID, Remark: "yzy3"}
	//db.Create(userFriends2)
	//
	//db.Model(user1).Association("Friends").Delete(user2)
	//db.Model(user1).Association("Friends").Append(user2)
}

func testGroup() {
	groupService := service.GroupService{}
	//db := database.DB

	//owner, _ := service.GetUserByName("yzy")
	//group := &model.Group{Name: "group", Owner: owner}
	//db.Create(group)
	//user1, _ := service.GetUserByName("y")
	//user2, _ := service.GetUserByName("yy")
	//db.Model(group).Association("Members").Append(user1)
	//db.Model(group).Association("Members").Append(user2)

	member, _ := groupService.GetGroupMemberNameList(2)
	log.Infof("%+v", member)
}

func main() {
	// testRedis()
	// testProto()
	//testLogin()
	//testChat()
	//testDatabase()
	testGroup()
}
