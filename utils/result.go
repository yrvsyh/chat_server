package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type (
	Result struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data,omitempty"`
	}
)

func Success(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Result{Code: 0, Msg: msg})
}

func SuccessWithData(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Result{Code: 0, Msg: msg, Data: data})
}

func Error(c *gin.Context, code int, msg string) {
	ret := Result{Code: code, Msg: msg}
	log.Error(ret)
	c.JSON(http.StatusOK, ret)
}

func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	ret := Result{Code: code, Msg: msg, Data: data}
	log.Error(ret)
	c.JSON(http.StatusOK, ret)
}
