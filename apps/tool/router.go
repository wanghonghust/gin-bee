package tool

import (
	"gin-bee/apps/tool/api"
	"gin-bee/apps/tool/api/core"
	"gin-bee/middleware"

	"github.com/gin-gonic/gin"
)

func RouterHandler(r *gin.RouterGroup) {
	tGroup := r.Group("/tool")
	tGroup.GET("/ws/:id", core.WsSsh)
	tGroup.Use(middleware.Authenticate())
	{
		tGroup.POST("qr-code", api.CQRCodeController.GenerateQRCode)

		tGroup.GET("/system-info", api.CSystem.List)

		tGroup.POST("/async_task", api.CTask.Create)

		tGroup.GET("/async_task", api.CTask.List)

		tGroup.POST("/ssh/config", core.GetSshConfig)
	}

}
