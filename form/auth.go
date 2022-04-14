package form

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthForm struct {
	Username  string                `form:"username" binding:"required"`
	Password  string                `form:"password" binding:"required"`
	PublicKey *multipart.FileHeader `form:"public_key" binding:"required"`
}

func (form *AuthForm) isValid(c *gin.Context) bool {
	err := c.ShouldBind(form)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{})
	}

	return true
}
