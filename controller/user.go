package controller

import (
	"chat_server/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetUserAvatar(c *gin.Context) {
	fmt.Println(c.Param("id"))
	utils.Success(c, "avatar")
}
