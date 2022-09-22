package core

import (
	"bytes"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var sshCfg = make(chan SSHConfig, 1)

// GetSshConfig
// @Summary
// @Schemes
// @Description ssh配置
// @Tags
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /tool/ssh/config [post]
func GetSshConfig(c *gin.Context) {
	sCfg := SSHConfig{}
	err := c.Bind(&sCfg)
	sshCfg <- sCfg
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "ssh配置错误"})
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "配置成功"})
}

// WsSsh
// @Summary
// @Schemes
// @Description shh连接
// @Tags
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /tool/ssh [get]
func WsSsh(c *gin.Context) {
	sCfg := <-sshCfg
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if handleError(c, err) {
		zaplog.Logger.Errorf("ws handdle error:%v", err)
		return
	}
	defer wsConn.Close()

	cols, err := strconv.Atoi(c.DefaultQuery("cols", "120"))
	if wshandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "32"))
	if wshandleError(wsConn, err) {
		return
	}
	client, err := NewSshClient(sCfg)
	if wshandleError(wsConn, err) {
		return
	}
	defer client.Close()
	ssConn, err := NewSshConn(cols, rows, client)
	if wshandleError(wsConn, err) {
		return
	}
	defer ssConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan

	zaplog.Logger.Info("websocket finished")
}
