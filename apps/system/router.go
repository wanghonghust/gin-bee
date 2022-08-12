package system

import (
	"github.com/gin-gonic/gin"
)

func RouterHandler(r *gin.Engine) {
	group := r.Group("/system")
	group.GET("/file", cSystem.File)
	//group.Use(middleware.Autenticate())
	group.POST("/fileupload", cSystem.FileUpload)

}
