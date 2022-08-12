package middleware

import (
	"gin-bee/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Autenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 1,
				"msg":  "未登录状态！",
			})
			c.Abort()
			return
		}
		if strings.HasPrefix(authorization, "Bearer ") {
			resJwt := strings.Split(authorization, " ")
			token := resJwt[len(resJwt)-1]
			data, err := utils.ParseToken(token)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": err.Error(),
				})
				c.Abort()
				return
			}
			if !data.State {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg": "用户被禁用",
				})
				c.Abort()
				return
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "token格式不正确",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
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
