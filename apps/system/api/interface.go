package api

import (
	"gin-bee/apps"
	"gin-bee/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

var CSystemInterface SystemInterfaceController

type SystemInterfaceController struct {
}

// GetApi
// @Summary
// @Schemes
// @Description 获取所有的api详情
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.APIInfos
// @Failure 400 {object} response.Response
// @Router /api/system/api [get]
func (sc *SystemInterfaceController) GetApi(c *gin.Context) {
	var apiInfos response.APIInfos
	if err := apps.Db.Find(&apiInfos.Data).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "查询数据失败"})
		return
	}
	c.JSONP(http.StatusOK, apiInfos)
}
