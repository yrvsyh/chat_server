package controller

import (
	"chat_server/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (AuthController) Register(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "auth/register.html", nil)
	} else if c.Request.Method == http.MethodPost {
		json := struct {
			Username  string `json:"username" binding:"required"`
			Password  string `json:"password" binding:"required"`
			PublicKey string `json:"public_key" binding:"required"`
		}{}

		if err := c.ShouldBind(&json); err != nil {
			Err(c, err)
			return
		}

		if err := userService.Register(json.Username, json.Password, json.PublicKey); err != nil {
			Err(c, err)
			return
		}

		Success(c)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func (AuthController) Login(c *gin.Context) {
	json := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	if err := c.ShouldBind(&json); err != nil {
		Err(c, err)
		return
	}

	id, username, publicKey, err := userService.Login(json.Username, json.Password)
	if err != nil {
		Err(c, err)
		return
	}

	session := middleware.GetAuthSession(c)
	session.Values["id"] = id
	session.Values["username"] = username

	if err := session.Save(c.Request, c.Writer); err != nil {
		Err(c, err)
		return
	}

	SuccessData(c, gin.H{"id": id, "public_key": publicKey})
}

func (AuthController) Logout(c *gin.Context) {
	session := middleware.GetAuthSession(c)
	delete(session.Values, "id")
	delete(session.Values, "username")
	if err := session.Save(c.Request, c.Writer); err != nil {
		Err(c, err)
		return
	}

	Success(c)
}

func GetLoginUserInfo(c *gin.Context) (uint32, string) {
	session := middleware.GetAuthSession(c)
	return session.Values["id"].(uint32), session.Values["username"].(string)
}
