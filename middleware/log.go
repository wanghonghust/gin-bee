package middleware

import (
	"bytes"
	"encoding/json"
	"gin-bee/apps"
	system "gin-bee/apps/system/model"
	"gin-bee/utils"
	"gin-bee/utils/constant"
	"gin-bee/zaplog"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogMiddleware() gin.HandlerFunc {

	// 请求日志记录
	var methodColorMap = map[string]int{
		"GET":    constant.BGBLUE,
		"DELETE": constant.BGRED,
		"POST":   constant.BGGREEN,
		"PUT":    constant.BGYELLOW,
	}
	var getMethodBg = func(method string) int {
		if val, ok := methodColorMap[method]; ok {
			return val
		} else {
			return 0
		}
	}
	var getStatusBg = func(status int) int {
		if 100 <= status && status < 200 {
			return constant.BGWHITE
		} else if 200 <= status && status < 300 {
			return constant.BGGREEN
		} else if 300 <= status && status < 400 {
			return constant.BGYELLOW
		} else {
			return constant.BGRED
		}
	}
	return func(c *gin.Context) {
		if c.ContentType() == "multipart/form-data" {
			return
		}
		// 拷贝请求体
		reqBody, _ := ioutil.ReadAll(c.Request.Body)
		var m map[string]interface{}
		newBody := make([]byte, 0)
		// 当请求体长度为0时,解析会出错
		if len(reqBody) > 0 {
			err := json.Unmarshal(reqBody, &m)
			if err != nil {
				zaplog.Logger.Errorf("Unmarshal reqBody Error:%v", err)
				return
			}
			newBody, err = json.Marshal(m)
			if err != nil {
				zaplog.Logger.Errorf("Marshal newBody Error:%v", err)
				return
			}
		}
		// 将原始的body重新赋值给context
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(newBody))
		start := time.Now()
		var blw bodyLogWriter
		blw = bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()
		// 获取响应
		strBody := ""
		strBody = strings.Trim(blw.bodyBuf.String(), " ")

		url := c.Request.URL.String()
		method := c.Request.Method
		cost := time.Since(start)
		userId, _ := c.Get("user_id")
		uid, _ := utils.AnyToUintPtr(userId)
		zaplog.Logger.Infof("\u001B[%d;37m%-6s\u001B[0m    \u001B[%d;37m%3d\u001B[0m   path:%s    body:%s  time:%v    user_id:%v", getMethodBg(method), method, getStatusBg(c.Writer.Status()), c.Writer.Status(), url, reqBody, cost, userId)
		// 保存到数据库
		if method != "GET" {
			log := system.Log{UserId: &uid, RemoteIP: c.RemoteIP(), ResponseTime: float64(uint(cost.Nanoseconds())) / 1e6, FullPath: c.FullPath(), Method: method, Body: string(reqBody), Response: strBody, Status: c.Writer.Status()}
			if err := apps.Db.Create(&log).Error; err != nil {
				zaplog.Logger.Error(err)
				c.Abort()
				return
			}
		}
	}
}
