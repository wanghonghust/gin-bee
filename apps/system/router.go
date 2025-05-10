package system

import (
	"gin-bee/apps/system/api"
	"gin-bee/middleware"

	"github.com/gin-gonic/gin"
)

func RouterHandler(r *gin.RouterGroup) {
	sysGroup := r.Group("/system")
	sysGroup.GET("/file", api.CSystem.File)
	sysGroup.GET("/log", api.CLog.Logs)
	sysGroup.Use(middleware.Authenticate())
	{
		sysGroup.POST("/file", api.CSystem.FileUpload)

		sysGroup.GET("/menu", api.CSystem.Menus)
		sysGroup.POST("/menu", api.CSystem.AddMenu)
		sysGroup.PUT("/menu", api.CSystem.EditMenu)
		sysGroup.DELETE("/menu", api.CSystem.DeleteMenu)
		sysGroup.GET("/menu/test", api.CSystem.TestMenu)

		sysGroup.GET("/role", api.CRole.Roles)
		sysGroup.POST("/role", api.CRole.AddRole)
		sysGroup.PUT("/role", api.CRole.EditRole)
		sysGroup.DELETE("/role", api.CRole.DeleteRole)

		sysGroup.GET("/api", api.CSystemInterface.GetApi)
	}

	permGroup := r.Group("/system")
	permGroup.Use(middleware.Authenticate())
	{
		permGroup.GET("/perm", api.CPermission.Permissions)
		permGroup.POST("/perm", api.CPermission.AddPermission)
		permGroup.PUT("/perm", api.CPermission.EditPermission)
		permGroup.DELETE("/perm", api.CPermission.DeletePermission)
	}
	group := r.Group("/auth")
	group.POST("/login", api.CAuth.Login)
	group.POST("", api.CAuth.Auth)
	group.Use(middleware.Authenticate())
	{
		group.POST("/user/create", api.CAuth.CreateUser)
		group.GET("/user", api.CAuth.UserInfo)
		group.PUT("/user/avatar", api.CAuth.EditUserAvatar)
		group.PUT("/user/limiter", api.CAuth.EditUserLimiter)
		group.PUT("/passwd", api.CAuth.ChangePwd)
		group.POST("/user", api.CAuth.AllUser)
		group.PUT("/user", api.CAuth.UpdateUserInfo)
		group.DELETE("/user", api.CAuth.DeleteUSer)
	}
}
