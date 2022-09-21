package core

import (
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"time"
)

func jsonError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": msg})
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		zaplog.Logger.Errorf("gin context http handler error:%v", err)
		jsonError(c, err.Error())
		return true
	}
	return false
}

func wshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		zaplog.Logger.Errorf("handler ws error:%v", err)
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			zaplog.Logger.Errorf("websocket writes control message failed:%v", err)
		}
		return true
	}
	return false
}
