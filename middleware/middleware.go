package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin-bee/apps"
	system "gin-bee/apps/system/model"
	"gin-bee/utils"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func Autenticate() gin.HandlerFunc {
	// 权限验证中间件
	return func(c *gin.Context) {
		// 登录验证
		userdata, err := GetCurrentUser(c)
		if err != nil {
			c.JSONP(http.StatusUnauthorized, gin.H{"msg": err.Error()})
			c.Abort()
			return
		}
		var user system.User
		user.ID = userdata.Id
		// 接口权限验证
		err = PathPermission(c, user)
		if err != nil {
			c.JSONP(http.StatusForbidden, gin.H{"msg": err.Error()})
			c.Abort()
			return
		}
		c.Set("user_id", user.ID)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	// 跨域中间件
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length,token,Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Headers")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}

func AccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//zaplog.Logger.Info(c.Request.RemoteAddr)

		if c.Request.Method == "POST" {
			var body map[string]any
			err := c.ShouldBindBodyWith(&body, binding.JSON)
			if err != nil {
				c.Next()
			}
			marshal, err := json.Marshal(body)
			if err != nil {
				fmt.Println(err)
			}
			zaplog.Logger.Info(string(marshal))
			zaplog.Logger.Info(c.RemoteIP())
			zaplog.Logger.Info(c.Request.Method)
		}
		c.Next()
	}
}

func PathPermission(c *gin.Context, usr system.User) (err error) {
	var perm system.Permission
	method := c.Request.Method
	path := c.FullPath()
	permName := fmt.Sprintf("%s:%s", method, path)
	err = apps.Db.Select("permissions.id").Joins("left join role_permissions on permissions.id = role_permissions.permission_id ").
		Joins("left join user_roles on role_permissions.role_id = user_roles.role_id").
		Joins("left join users on user_roles.user_id = users.id").
		Where("users.id = ?", usr.ID).
		Where("permissions.name = ? and permissions.type = 'path'", permName).
		Take(&perm).Error
	if err == gorm.ErrRecordNotFound {
		err = errors.New("用户无访问权限")
		return
	}
	return nil
}

func GetCurrentUser(c *gin.Context) (data *utils.JwtClaims, err error) {
	// 获取当前用户
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		err = errors.New("用户未登录")
		return nil, err
	} else if strings.HasPrefix(authorization, "Bearer ") {
		resJwt := strings.Split(authorization, " ")
		token := resJwt[len(resJwt)-1]
		data, err1 := utils.ParseToken(token)

		if err1 != nil {
			err = errors.New("token格式不正确")
			return nil, err
		}

		if !data.State {
			err = errors.New("用户已被停用")
			return nil, err
		}
		return data, nil
	} else {
		err = errors.New("token格式不正确")
		return nil, err
	}

}
