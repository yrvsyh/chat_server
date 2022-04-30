package controller

import (
	e "chat_server/errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context) {
	c.Set("Code", 0)
	c.Set("Msg", "success")
	c.JSON(http.StatusOK, Result{Code: 0, Msg: "success"})
}

func SuccessData(c *gin.Context, data interface{}) {
	c.Set("Code", 0)
	c.Set("Msg", "success")
	c.JSON(http.StatusOK, Result{Code: 0, Msg: "success", Data: data})
}

func Err(c *gin.Context, err error) {
	if err != nil {
		logrus.Warn(err.Error())
		msg := fmt.Sprintf("%s", err)
		c.Set("Code", -1)
		c.Set("Msg", msg)
		c.JSON(http.StatusOK, Result{Code: -1, Msg: msg})
		c.Abort()
	}
}

func Error(c *gin.Context, err error, msg string) {
	err = e.Wrap(err, msg)
	Err(c, err)
}

func ErrorCode(c *gin.Context, code int, msg string) {
	c.Set("Code", code)
	c.Set("Msg", msg)
	c.JSON(http.StatusOK, Result{Code: code, Msg: msg})
	c.Abort()
}

func ErrorData(c *gin.Context, code int, msg string, data interface{}) {
	c.Set("Code", code)
	c.Set("Msg", msg)
	c.JSON(http.StatusOK, Result{Code: code, Msg: msg, Data: data})
	c.Abort()
}
