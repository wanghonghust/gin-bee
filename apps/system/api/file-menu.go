package api

import (
	"fmt"
	"gin-bee/apps"
	"gin-bee/apps/system/model"
	"gin-bee/apps/system/request"
	"gin-bee/config"
	"gin-bee/response"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	CSystem = System{}
)

type System struct {
}

// File
// @Summary
// @Schemes
// @Description 下载文件
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce image/png,image/gif,image/jpeg,application/octet-stream
// @Param id query int true "文件id" mininum(1) maxinum(100)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/file [get]
func (s *System) File(c *gin.Context) {
	id := c.Query("id")
	parsedId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求参数不正确",
		})
		return
	}
	file := model.File{}
	file.Where("id = ?", parsedId).First(&file)
	if file.Path != "" {
		c.File(file.Path)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "文件不存在",
		})
	}

}

// FileUpload
// @Summary
// @Schemes
// @Description 上传文件
// @Tags
// @Security ApiKeyAuth
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/file [post]
func (s *System) FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "上传失败"})
		return
	}
	filename := file.Filename
	uploadPath := config.Cfg.Upload.File
	err = os.MkdirAll(uploadPath, os.FileMode(777))
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "上传失败"})
		return
	}
	id, _ := gonanoid.New(20)
	filename = id + "_" + filename
	dst := path.Join(uploadPath, filename)
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "上传失败"})
		return
	}

	var uploadFile model.File
	uploadFile.Path = dst
	tx := apps.Db.Create(&uploadFile)
	if tx.Error != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "上传失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "上传成功", "id": uploadFile.ID})
}

// Menus
// @Summary
// @Schemes
// @Description 获取树形结构菜单
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} response.MenuResponse
// @Failure 400 {object} response.Response
// @Router /api/system/menu [get]
func (s *System) Menus(c *gin.Context) {
	menus, err := getMenu(nil, 0, 0)
	if err != nil {
		c.JSONP(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"menus": menus})
}

// AddMenu
// @Summary
// @Schemes
// @Description 新增菜单
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.MenuAddParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/menu [post]
func (s *System) AddMenu(c *gin.Context) {
	var param request.MenuAddParam
	var menu model.Menu

	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		response.BadRequest(c)
		return
	}
	if err = param.Validator(); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	err = c.ShouldBindBodyWith(&menu, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	fmt.Println(menu)
	if err = apps.Db.Create(&menu).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "新增成功"})
}

// EditMenu
// @Summary
// @Schemes
// @Description 编辑菜单
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param object body request.MenuEditParam true "请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/menu [put]
func (s *System) EditMenu(c *gin.Context) {
	var param request.MenuEditParam
	var menu model.Menu

	err := c.ShouldBindBodyWith(&param, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if err = param.Validator(); err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	err = c.ShouldBindBodyWith(&menu, binding.JSON)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if err = apps.Db.Select("parent_id", "label", "icon", "link", "permission_sign", "local").Updates(&menu).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "修改成功"})
}

// DeleteMenu
// @Summary
// @Schemes
// @Description 删除菜单
// @Tags
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id body int true "int valid" mininum(1)
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/system/menu [delete]
func (s *System) DeleteMenu(c *gin.Context) {
	var menu model.Menu
	err := c.BindJSON(&menu)
	if err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	if menu.ID == 0 {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "请填写参数：id"})
		return
	}
	// 删除时也删除所有关联
	if err = apps.Db.Unscoped().Select(clause.Associations).Delete(&menu).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"msg": "删除成功"})
}

// 递归获取所有菜单
func getMenu(pid *uint, page int, pageSize int) ([]response.TreeMenu, error) {
	// num 为一级菜单返回的数量，为0返回所有
	var menus []model.Menu
	var tx *gorm.DB

	if pid == nil {
		// 查询pid为空的情况，使用原生sql，gorm无法查询
		if page != 0 {
			tx = apps.Db.Raw(fmt.Sprintf("SELECT * FROM `menus` WHERE `menus`.`parent_id` is null AND `menus`.`deleted_at` IS NULL LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)).Scan(&menus)
		} else {
			tx = apps.Db.Raw("SELECT * FROM `menus` WHERE `menus`.`parent_id` is null AND `menus`.`deleted_at` IS NULL").Scan(&menus)
		}

	} else {
		tx = apps.Db.Where(&model.Menu{ParentId: pid}).Order("id").Find(&menus)
	}

	var treeMenu []response.TreeMenu
	for _, menu := range menus {
		child, err := getMenu(&menu.ID, 0, 0)
		if err != nil {
			return make([]response.TreeMenu, 0), err
		}
		node := response.TreeMenu{
			ID:       menu.ID,
			Label:    menu.Label,
			Local:    menu.Local,
			Link:     menu.Link,
			Icon:     menu.Icon,
			ParentId: menu.ParentId,
			CreateAt: menu.CreatedAt.ToString(),
			Children: child,
		}
		treeMenu = append(treeMenu, node)
	}
	return treeMenu, tx.Error
}

func (s *System) TestMenu(c *gin.Context) {
	var menus []model.Menu
	apps.Db.Where("id in (?)", []uint{2, 4, 5, 51, 59}).Find(&menus)
	var ids []uint
	var ids1 []uint
	for _, item := range menus {
		ids1 = append(ids1, item.ID)
	}
	for _, item := range ids1 {
		getTreeId(item, &ids)
	}
	roleMenus, _ := TreeOfMenus(menus)
	c.JSONP(http.StatusOK, gin.H{"msg": "查询成功", "menus": roleMenus})
}

func getTreeMenu(pid *uint, ids []uint) ([]response.TreeMenu, error) {
	// 根据id查询出树形菜单
	var menus []model.Menu
	var tx *gorm.DB

	if pid == nil {
		// 查询pid为空的情况，使用原生sql，gorm无法查询
		tx = apps.Db.Raw("SELECT * FROM `menus` WHERE `menus`.`parent_id` is null AND `menus`.`deleted_at` IS NULL AND id IN (?)", ids).Scan(&menus)
	} else {
		tx = apps.Db.Where(&model.Menu{ParentId: pid}).Where("id in (?)", ids).Order("id").Find(&menus)
	}

	var treeMenu []response.TreeMenu
	for _, menu := range menus {
		child, err := getTreeMenu(&menu.ID, ids)
		if err != nil {
			return make([]response.TreeMenu, 0), err
		}
		node := response.TreeMenu{
			ID:       menu.ID,
			Label:    menu.Label,
			Local:    menu.Local,
			Link:     menu.Link,
			Icon:     menu.Icon,
			ParentId: menu.ParentId,
			CreateAt: menu.CreatedAt.ToString(),
			Children: child,
		}
		treeMenu = append(treeMenu, node)
	}
	return treeMenu, tx.Error
}

func isIn(nums []uint, num uint) (status bool) {
	// 判断数字是否在数组中
	for _, item := range nums {
		if item == num {
			return true
		}
	}
	return false
}

func getTreeId(id uint, ids *[]uint) {
	// 递归查找树形结构的id
	var menu model.Menu
	if isIn(*ids, id) {
		return
	}
	*ids = append(*ids, id)
	apps.Db.Where("id = ?", id).Find(&menu)
	if menu.ParentId != nil {
		getTreeId(*menu.ParentId, ids)
	}
}

func TreeOfMenus(menus []model.Menu) (treeMenu []response.TreeMenu, err error) {
	// 存储整个树形结构的所有menu.id
	var ids []uint
	// 存储menus 的 menu.id
	var menuId []uint

	for _, item := range menus {
		menuId = append(menuId, item.ID)
	}
	for _, item := range menuId {
		getTreeId(item, &ids)
	}
	treeMenu, err = getTreeMenu(nil, ids)
	return
}
