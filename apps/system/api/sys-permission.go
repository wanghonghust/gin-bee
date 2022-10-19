package api

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	"gin-bee/apps/system/model"
	"gin-bee/apps/system/request"
	"gin-bee/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm/clause"
	"net/http"
)

var (
	CPermission = PermissionController{}
)

type PermissionController struct {
}

// Permissions
// @Summary
// @Schemes
// @Description 获取所有权限
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/perm [get]
func (pc *PermissionController) Permissions(c *gin.Context) {
	var perms []model.Permission
	// 用于时间格式化输出
	var res []map[string]any
	if err := apps.Db.Find(&perms).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据查询失败"})
		return
	}
	marshal, err := json.Marshal(&perms)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据Marshal失败"})
		return
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据Unmarshal失败"})
		return
	}
	// 时间格式化C
	for index, item := range res {
		var pid uint
		var count int64
		switch item["id"].(type) {
		case float64:
			pid = uint(item["id"].(float64))
			break
		}
		res[index]["createdAt"] = utils.StrTimeFormat(fmt.Sprintf("%v", item["createdAt"]))
		// 检查权限是否被引用
		apps.Db.Table("role_permissions").Where("permission_id = ?", pid).Count(&count)
		res[index]["deleteAble"] = count == 0
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "请求成功", "data": res})
}

// AddPermission
// @Summary
// @Schemes
// @Description 新增权限
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.AddPermissionParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/perm [post]
func (pc *PermissionController) AddPermission(c *gin.Context) {
	var param request.AddPermissionParam
	var perm model.Permission
	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}
	fmt.Println(param)
	if err = param.Validator(); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	err = c.ShouldBindBodyWith(&perm, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}
	if err = apps.Db.Create(&perm).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "新建成功"})
}

// EditPermission
// @Summary
// @Schemes
// @Description 编辑权限
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.EditPermissionParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/perm [put]
func (pc *PermissionController) EditPermission(c *gin.Context) {
	var param request.EditPermissionParam
	var perm model.Permission
	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}
	if err = param.Validator(); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	err = c.ShouldBindBodyWith(&perm, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}
	if err = apps.Db.Select("name", "type", "desc").Updates(&perm).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "修改成功"})
}

// DeletePermission
// @Summary
// @Schemes
// @Description 删除权限
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.DeletePermissionParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/perm [delete]
func (pc *PermissionController) DeletePermission(c *gin.Context) {
	// 批量删除
	var param request.DeletePermissionParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请求参数不正确"})
		return
	}
	perms := make([]model.Permission, len(param.Id))
	for _, item := range param.Id {
		perms = append(perms, model.Permission{Model: apps.Model{ID: item}})
	}
	// 删除时也删除关联
	if err = apps.Db.Debug().Unscoped().Select(clause.Associations).Delete(&perms).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "删除成功"})
}
