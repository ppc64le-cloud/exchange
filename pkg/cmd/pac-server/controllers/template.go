package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server/internal"
)

var TemplateControllers = &template{}

type template struct {
}

func (u *template) GetTemplate(c *gin.Context) {
	name := c.Param("name")
	internal.Config.Backend.GetTemplate(c, name)
}

func (u *template) ListTemplate(c *gin.Context) {

}
