package api

import (
	"gin-bee/apps"
	"gin-bee/apps/system/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	CLog = LogController{}
)

type LogController struct {
}

// Logs
// @Summary
// @Schemes
// @Description 请求日志
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.LogResponse
// @Failure 400 {object} response.Response
// @Router /api/system/log [get]
func (lc *LogController) Logs(c *gin.Context) {
	var logs []model.Log
	if err := apps.Db.Find(&logs).Error; err != nil {
		c.JSONP(http.StatusNotFound, gin.H{"msg": "查询数据出错"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"data": logs})
}
