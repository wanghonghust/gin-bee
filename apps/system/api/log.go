package api

import (
	"gin-bee/apps"
	"gin-bee/apps/system/model"
	"github.com/asaskevich/govalidator"
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
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.LogResponse
// @Failure 400 {object} response.Response
// @Router /api/system/log [get]
func (lc *LogController) Logs(c *gin.Context) {
	var logs []model.Log
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "20")
	page64, _ := govalidator.ToInt(page)
	pageSize64, _ := govalidator.ToInt(pageSize)
	if page64 <= 0 {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "page must greater than 0!"})
		return
	}
	if pageSize64 <= 0 {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "pageSize must greater than 0!"})
		return
	}
	// 数据分页
	var total int64 // 总数据量
	if err := apps.Db.Model(model.Log{}).Count(&total).Order("id desc").Limit(int(pageSize64)).Offset(int((page64 - 1) * pageSize64)).Find(&logs).Error; err != nil {
		c.JSONP(http.StatusNotFound, gin.H{"msg": "查询数据出错"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"page": page64, "pageSize": pageSize64, "count": len(logs), "total": total, "data": logs})
}
