package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server/controllers"
)

func templatesRoutes(superRoute *gin.RouterGroup) {
	templatesRouter := superRoute.Group("/templates")
	{
		templatesRouter.GET("/", controllers.TemplateControllers.ListTemplate)
	}
	templateRouter := superRoute.Group("/template")
	{
		templateRouter.GET("/:name", controllers.TemplateControllers.GetTemplate)
	}
}
