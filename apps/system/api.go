package system

import (
	"gin-bee/apps"
	"gin-bee/config"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"net/http"
	"os"
	"path"
)

var (
	cSystem = System{}
)

type System struct {
}

func (s *System) File(c *gin.Context) {
	id := c.Query("id")
	file := File{}
	file.Where("id = ?", id).First(&file)
	if file.Path != "" {
		c.File(file.Path)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "文件不存在",
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

	var uploadFile File
	uploadFile.Path = dst
	tx := apps.Db.Create(&uploadFile)
	if tx.Error != nil {
		c.JSONP(http.StatusBadRequest, gin.H{"msg": "上传失败"})
		return
	}
	c.JSONP(http.StatusOK, gin.H{"msg": "上传成功", "id": uploadFile.ID})
}
