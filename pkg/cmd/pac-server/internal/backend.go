package internal

import "github.com/gin-gonic/gin"

type Backend interface {
	GetTemplate(c *gin.Context, name string)
}
