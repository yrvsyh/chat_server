package form

import "github.com/gin-gonic/gin"

type Form interface {
	isValid(c *gin.Context) bool
}
