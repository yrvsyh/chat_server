package middleware

import (
	"chat_server/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(config.SessionKey))

func SessionAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := GetAuthSession(c)
		_, ok := session.Values[config.SessionUserKey]
		if !ok {
			c.Status(http.StatusUnauthorized)
			c.Abort()
		}
	}
}

func GetAuthSession(c *gin.Context) *sessions.Session {
	session, _ := store.Get(c.Request, config.CookieKey)
	return session
}
