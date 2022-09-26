package middleware

import (
	"bytes"
	"encoding/json"
	"gin-bee/apps"
	system "gin-bee/apps/system/model"
	"gin-bee/utils"
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
		userId, ok := c.Get("user_id")
		if ok != true {
			c.Abort()
			return
		}
		uid, err := utils.AnyToUintPtr(userId)
		if err != nil {
			zaplog.Logger.Errorf("Convert Error:%v", err)
			c.Abort()
			return
		}
		// 保存到数据库
		if method != "GET" {
			log := system.Log{UserId: &uid, RemoteIP: c.RemoteIP(), ResponseTime: float64(uint(cost.Nanoseconds())) / 1e6, FullPath: c.FullPath(), Method: method, Body: string(reqBody), Response: strBody, Status: c.Writer.Status()}
			if err = apps.Db.Create(&log).Error; err != nil {
				zaplog.Logger.Error(err)
				c.Abort()
				return
			}
		}
		zaplog.Logger.Infof("[GIN] path:%s    method:%s	status:%d   body:%s  time:%v    user_id:%v", url, method, c.Writer.Status(), reqBody, cost, userId)
	}
}
