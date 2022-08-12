package auth

import (
	"gin-bee/middleware"
	"github.com/gin-gonic/gin"
)

func RouterHandler(r *gin.Engine) {
	group := r.Group("/auth")
	group.POST("/users/create", cAuth.CreateUser)
	group.POST("/login", cAuth.Login)
	group.POST("", cAuth.Auth)
	group.Use(middleware.Autenticate())
	group.POST("/userInfo", cAuth.UserInfo)
	group.POST("/user/avatar/edit", cAuth.EditUserAvatar)
	group.POST("/passwd/change", cAuth.ChangePwd)
	group.POST("/users", cAuth.AllUser)
	group.POST("/users/update", cAuth.UpdateUserInfo)
	group.POST("/users/delete", cAuth.DeleteUSer)
}
