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

func Success(c *gin.Context, msg string) {
	c.Set("Code", 0)
	c.Set("Msg", msg)
	c.JSON(http.StatusOK, Result{Code: 0, Msg: msg})
}

func SuccessWithData(c *gin.Context, msg string, data interface{}) {
	c.Set("Code", 0)
	c.Set("Msg", msg)
	c.JSON(http.StatusOK, Result{Code: 0, Msg: msg, Data: data})
}

func Error(c *gin.Context, code int, msg string) {
	c.Set("Code", code)
	c.Set("Msg", msg)
	ret := Result{Code: code, Msg: msg}
	c.JSON(http.StatusOK, ret)
}

func ErrorWithData(c *gin.Context, code int, msg string, data interface{}) {
	c.Set("Code", code)
	c.Set("Msg", msg)
	ret := Result{Code: code, Msg: msg, Data: data}
	c.JSON(http.StatusOK, ret)
}
