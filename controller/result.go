package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context) {
	c.Set("Code", 0)
	c.Set("Msg", "success")
	c.JSON(http.StatusOK, Result{Code: 0})
}

func SuccessData(c *gin.Context, data interface{}) {
	c.Set("Code", 0)
	c.Set("Msg", "success")
	c.JSON(http.StatusOK, Result{Code: 0, Data: data})
}

func Err(c *gin.Context, err error) {
	c.Set("Code", -1)
	c.Set("Msg", err.Error())
	c.JSON(http.StatusOK, Result{Code: -1, Msg: err.Error()})
	c.Abort()
}

func Error(c *gin.Context, msg string) {
	c.Set("Code", -1)
	c.Set("Msg", msg)
	c.JSON(http.StatusOK, Result{Code: -1, Msg: msg})
	c.Abort()
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
