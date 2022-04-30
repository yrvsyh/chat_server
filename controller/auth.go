package controller

import (
	"chat_server/middleware"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (AuthController) Register(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "auth/register.html", nil)
	} else if c.Request.Method == http.MethodPost {
		data := struct {
			Username  string `form:"username" json:"username" binding:"required"`
			Password  string `form:"password" json:"password" binding:"required"`
			PublicKey string `form:"public_key" json:"public_key" binding:"required"`
		}{}

		if err := c.ShouldBind(&data); err != nil {
			Err(c, err)
			return
		}

		if err := userService.Register(data.Username, data.Password, data.PublicKey); err != nil {
			Err(c, err)
			return
		}

		Success(c)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func (AuthController) Login(c *gin.Context) {
	data := struct {
		Username string `form:"username" json:"username" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}{}

	if err := c.ShouldBind(&data); err != nil {
		Err(c, err)
		return
	}

	user, err := userService.Login(data.Username, data.Password)
	if err != nil {
		Err(c, err)
		return
	}

	session := middleware.GetAuthSession(c)
	session.Values["id"] = user.ID
	session.Values["username"] = user.Username

	if err := session.Save(c.Request, c.Writer); err != nil {
		Err(c, err)
		return
	}

	type VO struct {
		ID        uint32 `json:"id"`
		Name      string `json:"name"`
		PublicKey string `json:"public_key"`
		Avatar    string `json:"avatar"`
	}

	var ret = VO{
		ID:        user.ID,
		Name:      user.Username,
		PublicKey: user.PublicKey,
		Avatar:    user.Avatar,
	}
	if name := strings.TrimSpace(user.Nickname); name != "" {
		ret.Name = name
	}

	SuccessData(c, ret)
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
