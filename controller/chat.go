package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func checkToken(token string) bool {
	return true
}

func ChatHandle(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// mt, token, err := ws.ReadMessage()
	// if mt != websocket.TextMessage || err != nil {
	// 	return
	// }
	// if checkToken(string(token)) {
	for {
		mt, message, err := ws.ReadMessage()
		if websocket.IsUnexpectedCloseError(err) {
			log.Printf("read: %#v", err)
			break
		}
		log.Printf("recv: %s", message)
		err = ws.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
		log.Printf("done")
	}
	// }
}
