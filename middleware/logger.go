package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func LoggerMiddleware(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Set("Code", 0)
		c.Set("Msg", "")
		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.Error(e)
			}
		} else {
			field := log.Fields{
				"status":     c.Writer.Status(),
				"method":     c.Request.Method,
				"path":       path,
				"time":       end.Format("2006-01-02 15:04:05"),
				"user-agent": c.Request.UserAgent(),
				"latency":    fmt.Sprintf("%v", latency),
				// "timestamp":  end.Unix(),
				// "ip":         c.ClientIP(),
			}
			if query != "" {
				field["query"] = query
			}
			code, _ := c.Get("Code")
			msg, _ := c.Get("Msg")
			logger.WithFields(field).Infof("REQUEST [%s %s] [Code=%d Msg=%s]", c.Request.Method, path, code, msg)
		}
	}
}
