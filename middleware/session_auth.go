package middleware

import (
	"chat_server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("session_key"))

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := GetAuthSession(c)
		_, ok := session.Values["id"]
		_, ok = session.Values["username"]
		if !ok {
			utils.Error(c, -1, "用户未登录")
			c.Abort()
		}
		//c.Set("id", id)
		//c.Set("username", username)
	}
}

func GetAuthSession(c *gin.Context) *sessions.Session {
	session, _ := store.Get(c.Request, "auth_info")
	return session
}
