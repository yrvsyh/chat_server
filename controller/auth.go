package controller

import (
	"chat_server/middleware"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (AuthController) Register(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "auth/register.html", nil)
	} else if c.Request.Method == http.MethodPost {
		form := struct {
			Username  string                `form:"username" binding:"required"`
			Password  string                `form:"password" binding:"required"`
			PublicKey *multipart.FileHeader `form:"public_key" binding:"required"`
		}{}

		if err := c.ShouldBind(&form); err != nil {
			Err(c, err)
			return
		}

		keyFile, err := form.PublicKey.Open()
		if err != nil {
			Err(c, err)
			return
		}

		publicKey, err := io.ReadAll(keyFile)
		if err != nil {
			Err(c, err)
			return
		}

		if err := userService.Register(form.Username, form.Password, publicKey); err != nil {
			Err(c, err)
			return
		}

		Success(c)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func (AuthController) Login(c *gin.Context) {
	form := struct {
		Username  string                `form:"username" binding:"required"`
		Password  string                `form:"password" binding:"required"`
		PublicKey *multipart.FileHeader `form:"public_key" binding:"-"`
	}{}

	if err := c.ShouldBind(&form); err != nil {
		Err(c, err)
		return
	}

	// keyFile, err := data.PublicKey.Open()
	// if err != nil {
	// 	Err(c, err)
	// 	return
	// }

	// publicKey, err := io.ReadAll(keyFile)
	// if err != nil {
	// 	Err(c, err)
	// 	return
	// }
	var publicKey []byte

	id, username, err := userService.Login(form.Username, form.Password, publicKey)
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

	SuccessData(c, gin.H{"id": id})
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
