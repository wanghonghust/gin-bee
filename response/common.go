package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Msg string
}

func Success(c *gin.Context, msg string) {
	c.JSONP(http.StatusOK, gin.H{"msg": msg})
}

func BadRequest(c *gin.Context) {
	c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求错误"})
}

func UnAuthed(c *gin.Context) {
	c.JSONP(http.StatusUnauthorized, gin.H{"msg": "未授权"})
}

func ServerError(c *gin.Context) {
	c.JSONP(http.StatusInternalServerError, gin.H{"msg": "服务端错误"})
}
