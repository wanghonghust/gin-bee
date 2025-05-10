package middleware

import (
	"errors"
	"fmt"
	"gin-bee/apps"
	system "gin-bee/apps/system/model"
	"gin-bee/redis"
	"gin-bee/response"
	"gin-bee/utils"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	goredis "github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AuthErr struct {
	Code int    `json:"code"` // 1000:OK,1001:unauthorized,1002:token expired,1003:cannot parse token,1004:user status is forbidden,1005:not a token string.
	Msg  string `json:"msg"`
}

func Authenticate() gin.HandlerFunc {
	// 权限验证中间件
	return func(c *gin.Context) {
		// 登录验证
		userdata, err := GetCurrentUser(c)
		if err != nil {
			c.JSONP(http.StatusUnauthorized, *err)
			c.Abort()
			return
		}
		var user system.User
		user.ID = userdata.Id
		c.Set("user_id", user.ID)
		if userdata.IsSuperUser {
			c.Next()
			return
		}
		// 接口权限验证
		err1 := PathPermission(c, user)
		if err1 != nil {
			c.JSONP(http.StatusForbidden, gin.H{"msg": err1.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	// 跨域中间件
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}

func AccessLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := GetCurrentUser(c)
		if data == nil {
			// 匿名用户放行
			return
		}
		limitKey := fmt.Sprintf("rate:%d", data.Id)
		if data.Limiter.On {
			_, err := redis.RedisCli.Get(c, limitKey).Result()
			// 键值不存在，初始值设为1。
			if err == goredis.Nil {
				redis.RedisCli.Set(c, limitKey, 1, time.Hour*24)
			} else { // 键值存在则加1
				redis.RedisCli.Incr(c, limitKey)
			}
			count, err := govalidator.ToInt(redis.RedisCli.Get(c, limitKey).Val())
			if err == nil && count > int64(data.Limiter.Limit) {
				c.JSONP(http.StatusNotAcceptable, response.LimitError{
					Msg: fmt.Sprintf("访问次数限制上限为：%d！", 1000),
				})
				c.Abort()
				return
			}
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

func GetCurrentUser(c *gin.Context) (data *utils.JwtClaims, err *AuthErr) {
	// 获取当前用户
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		err = &AuthErr{Code: 1001, Msg: "用户未登录"}
		return nil, err
	} else if strings.HasPrefix(authorization, "Bearer ") {
		resJwt := strings.Split(authorization, " ")
		token := resJwt[len(resJwt)-1]
		data, err1 := utils.ParseToken(token)

		if err1 != nil {
			if strings.Contains(err1.Error(), "token is expired") {
				err = &AuthErr{Code: 1002, Msg: "token已过期"}
			} else {
				err = &AuthErr{Code: 1003, Msg: "无法解析token"}
			}

			return nil, err
		}
		var user system.User
		if err := apps.Db.Where("id = ?", data.Id).Find(&user).Error; err != nil {
			return nil, &AuthErr{Code: 1006, Msg: err.Error()}
		}
		// 从数据库更新数据
		data.Limiter = utils.Limiter{On: user.Limiter.On, Limit: user.Limiter.Limit}
		data.State = user.State

		if !data.State {
			err = &AuthErr{Code: 1004, Msg: "用户已被停用"}
			return nil, err
		}

		return data, nil

	} else {
		err = &AuthErr{Code: 1005, Msg: "不是一个正确的token"}
		return nil, err
	}

}
