package api

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	"gin-bee/apps/system/model"
	"gin-bee/apps/system/request"
	"gin-bee/response"
	"gin-bee/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
)

var (
	CRole = RoleController{}
)

type RoleController struct {
}

// Roles
// @Summary
// @Schemes
// @Description 获取所有角色
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.RoleResponse
// @Failure 400 {object} response.Response
// @Router /api/system/role [get]
func (r *RoleController) Roles(c *gin.Context) {
	var roles []model.Role
	var roleData []response.RoleData
	// 用于时间格式化输出
	var res []map[string]any
	// 使用预加载，查询出关联的权限
	if err := apps.Db.Preload("Permission").Preload("Menu").Find(&roles).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据查询失败"})
		return
	}
	// 获取树形菜单
	for _, item := range roles {
		var menuId []uint
		for _, me := range item.Menu {
			menuId = append(menuId, me.ID)
		}
		menu, err := TreeOfMenus(item.Menu)
		if menu == nil {
			menu = make([]response.TreeMenu, 0)
		}
		if err != nil {
			c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据查询失败"})
			return
		}
		tmpData := response.RoleData{Role: item, Menu: menu, MenuId: menuId}
		roleData = append(roleData, tmpData)
	}

	marshal, err := json.Marshal(&roleData)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据Marshal失败"})
		return
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "数据Unmarshal失败"})
		return
	}
	// 时间格式化
	for index, item := range res {
		res[index]["createdAt"] = utils.StrTimeFormat(fmt.Sprintf("%v", item["createdAt"]))
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "请求成功", "data": res})
}

// AddRole
// @Summary
// @Schemes
// @Description 新增角色
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.AddRoleParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/role [post]
func (r *RoleController) AddRole(c *gin.Context) {
	var role model.Role
	var param request.AddRoleParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	role.Name = param.Name
	// 将选择的权限id初始化为权限对象，并添加到角色权限中
	for _, item := range param.Permission {
		var tmpPerm model.Permission
		tmpPerm.Model = apps.Model{ID: item}
		role.Permission = append(role.Permission, tmpPerm)
	}

	// 将选择的菜单id初始化为菜单对象，并添加到角色菜单中
	for _, item := range param.Menu {
		var tmpMenu model.Menu
		tmpMenu.Model = apps.Model{ID: item}
		role.Menu = append(role.Menu, tmpMenu)
	}

	if err = apps.Db.Create(&role).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "新增成功"})
}

// EditRole
// @Summary
// @Schemes
// @Description 编辑角色
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.EditRoleParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/role [put]
func (r *RoleController) EditRole(c *gin.Context) {
	var role model.Role
	var param request.EditRoleParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	role.ID = param.Id
	role.Name = param.Name
	// 将选择的权限id初始化为权限对象，并添加到角色权限中
	for _, item := range param.Permission {
		var tmpPerm model.Permission
		tmpPerm.Model = apps.Model{ID: item}
		role.Permission = append(role.Permission, tmpPerm)
	}

	// 将选择的菜单id初始化为菜单对象，并添加到角色菜单中
	for _, item := range param.Menu {
		var tmpMenu model.Menu
		tmpMenu.Model = apps.Model{ID: item}
		role.Menu = append(role.Menu, tmpMenu)
	}

	// 更新自身字段
	if err = apps.Db.Updates(&role).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	// 更新关联关系
	if err = apps.Db.Model(&role).Association("Permission").Replace(role.Permission); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if err = apps.Db.Model(&role).Association("Menu").Replace(role.Menu); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"msg": "修改成功"})
}

// DeleteRole
// @Summary
// @Schemes
// @Description 删除角色
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.DeleteRoleParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/role [delete]
func (r *RoleController) DeleteRole(c *gin.Context) {
	var roles []model.Role
	var param request.DeleteRoleParam
	err := c.BindJSON(&param)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	for _, item := range param.Id {
		roles = append(roles, model.Role{Model: apps.Model{ID: item}})
	}
	// 开启事务
	err = apps.Db.Transaction(func(tx *gorm.DB) error {
		// 删除，彻底删除,并清除关联关系
		if err3 := tx.Unscoped().Select(clause.Associations).Delete(&roles).Error; err3 != nil {
			return err3
		}
		return nil
	})
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "删除成功"})
}
