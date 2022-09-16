package api

import (
	"fmt"
	"gin-bee/apps"
	"gin-bee/apps/system/model"
	"gin-bee/apps/system/request"
	"gin-bee/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path"
	"strconv"
)

var (
	CSystem = System{}
)

type TreeMenu struct {
	ID       uint        `json:"id"`
	Label    string      `json:"label"`
	ParentId *uint       `json:"parentId"`
	Link     string      `json:"link"`
	Icon     string      `json:"icon"`
	Children []*TreeMenu `json:"children"`
	CreateAt string      `json:"createAt"`
}
type System struct {
}

func (s *System) File(c *gin.Context) {
	id := c.Query("id")
	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请求参数不正确",
		})
		return
	}
	file := model.File{}
	file.Where("id = ?", id).First(&file)
	if file.Path != "" {
		c.File(file.Path)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "文件不存在",
		})
	}

}

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

func (s *System) Menus(c *gin.Context) {
	menus, err := getMenu(nil, 0, 0)
	if err != nil {
		c.JSONP(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"menus": menus})
}

func (s *System) AddMenu(c *gin.Context) {
	var param request.MenuAddParam
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
	if err = apps.Db.Create(&menu).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "新增成功"})
}

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
	if err = apps.Db.Debug().Select("parent_id", "label", "icon", "link").Updates(&menu).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "修改成功"})
}

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
	if err = apps.Db.Unscoped().Delete(&menu).Error; err != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	c.JSONP(http.StatusOK, gin.H{"msg": "删除成功"})
}

// 递归获取所有菜单
func getMenu(pid *uint, page int, pageSize int) ([]*TreeMenu, error) {
	// num 为一级菜单返回的数量，为0返回所有
	var menus []model.Menu
	var tx *gorm.DB

	if pid == nil {
		// 查询pid为空的情况，使用原生sql，gorm无法查询
		if page != 0 {
			tx = apps.Db.Debug().Raw(fmt.Sprintf("SELECT * FROM `menus` WHERE `menus`.`parent_id` is null AND `menus`.`deleted_at` IS NULL LIMIT %d OFFSET %d", pageSize, (page-1)*pageSize)).Scan(&menus)
		} else {
			tx = apps.Db.Raw("SELECT * FROM `menus` WHERE `menus`.`parent_id` is null AND `menus`.`deleted_at` IS NULL").Scan(&menus)
		}

	} else {
		tx = apps.Db.Where(&model.Menu{ParentId: pid}).Order("id").Find(&menus)
	}

	var treeMenu []*TreeMenu
	for _, menu := range menus {
		child, err := getMenu(&menu.ID, 0, 0)
		if err != nil {
			return nil, err
		}
		node := &TreeMenu{
			ID:       menu.ID,
			Label:    menu.Label,
			Link:     menu.Link,
			Icon:     menu.Icon,
			ParentId: menu.ParentId,
			CreateAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		node.Children = child
		treeMenu = append(treeMenu, node)
	}
	return treeMenu, tx.Error
}
