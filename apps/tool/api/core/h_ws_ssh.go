package core

import (
	"bytes"
	"fmt"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsSsh(c *gin.Context) {
	sshCfg := SSHConfig{User: "root", Password: "Emergency520", Addr: fmt.Sprintf("%s:%d", "121.4.61.20", 22)}
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
	client, err := NewSshClient(sshCfg)
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
