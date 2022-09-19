package model

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	"gin-bee/utils"
	"gin-bee/zaplog"
	"gorm.io/gorm/clause"
	"os"
	"strings"
)

type API struct {
	apps.Model
	Path        string `json:"path"  gorm:"type:varchar(64);required;not null;uniqueIndex:api;"`
	Method      string `json:"method" gorm:"type:varchar(64);required;not null;uniqueIndex:api;"`
	Description string `json:"description" gorm:"type:varchar(64);required;not null"`
}

func init() {
	var api API
	var perms []Permission
	err := api.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}
	zaplog.Logger.Info("数据表apis迁移成功")
	apis, err := getAllApiFromDoc()
	if err != nil {
		zaplog.Logger.Error("解析API文档失败")
		return
	}
	if err = apps.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&apis).Error; err != nil {
		zaplog.Logger.Error("初始化数据表apis数据失败")
		return
	}
	zaplog.Logger.Info("初始化数据表apis数据成功")

	// 初始化接口权限

	for _, item := range apis {
		perm := Permission{
			Type: "path",
			Name: fmt.Sprintf("%s:%s", item.Method, item.Path),
			Desc: item.Description,
		}
		perms = append(perms, perm)
	}
	if err = apps.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&perms).Error; err != nil {
		zaplog.Logger.Error("初始化接口权限失败")
		return
	}
	zaplog.Logger.Info("初始化接口权限成功")

}

func (a *API) Migrate() error {
	return apps.Db.AutoMigrate(&a)
}

func getAllApiFromDoc() (res []API, err error) {
	jsonPath := "./docs/swagger.json"
	var swagger utils.SwaggerJson
	file, err := os.ReadFile(jsonPath)
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &swagger)
	if err != nil {
		return nil, err
	}
	for path, val := range swagger.Paths {
		for method, detail := range val {
			var tmpApi API
			tmpApi.Path = path
			tmpApi.Method = strings.ToUpper(method)
			tmpApi.Description = detail["description"].(string)
			res = append(res, tmpApi)
		}
	}
	return
}
