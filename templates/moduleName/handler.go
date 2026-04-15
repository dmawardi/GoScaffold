package moduleName

import (
	"github.com/dmawardi/goTemplate/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, c ModuleNameController) {
	// /moduleNames
	moduleNames := rg.Group("/moduleNames")
	{
		// @tag.name Protected Routes
		// @tag.description Routes that require authentication
		moduleNames.Use(auth.AuthenticateJWT())
		moduleNames.GET("", c.FindAll)
		moduleNames.GET("/:id", c.Find)
		moduleNames.POST("", c.Create)
		moduleNames.PUT("/:id", c.Update)
		moduleNames.DELETE("/:id", c.Delete)
	}
}
