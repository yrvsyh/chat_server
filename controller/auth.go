package controller

import (
	"chat_server/middleware"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct{}

func (AuthController) getPublicKey(c *gin.Context) ([]byte, error) {
	var public_key []byte

	keyFormFile, err := c.FormFile("public_key")
	if err != nil {
		return public_key, err
	}

	keyFile, err := keyFormFile.Open()
	if err != nil {
		return public_key, err
	}

	public_key, err = ioutil.ReadAll(keyFile)
	return public_key, err
}

func (a AuthController) Register(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		c.HTML(http.StatusOK, "auth/register.html", nil)
	} else if c.Request.Method == http.MethodPost {
		data := &struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}{}

		if err := c.Bind(data); err != nil {
			Err(c, err)
			return
		}

		public_key, err := a.getPublicKey(c)
		if err != nil {
			Err(c, err)
			return
		}

		if err := userService.Register(data.Username, data.Password, public_key); err != nil {
			Err(c, err)
			return
		}

		Success(c)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func (a AuthController) Login(c *gin.Context) {
	data := &struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}{}

	if err := c.Bind(data); err != nil {
		Err(c, err)
		return
	}

	public_key, err := a.getPublicKey(c)
	// if err != nil {
	// 	Err(c, err)
	// 	return
	// }

	id, username, err := userService.Login(data.Username, data.Password, public_key)
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
