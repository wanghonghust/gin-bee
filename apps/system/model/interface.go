package model

import (
	"encoding/json"
	"fmt"
	"gin-bee/apps"
	"gin-bee/docs"
	"gin-bee/utils"
	"gin-bee/zaplog"
	"regexp"
	"strings"

	"gorm.io/gorm/clause"
)

type API struct {
	apps.Model
	Path        string `json:"path"  gorm:"type:varchar(64);required;not null;uniqueIndex:api;"`
	Method      string `json:"method" gorm:"type:varchar(64);required;not null;uniqueIndex:api;"`
	Description string `json:"description" gorm:"type:varchar(64);required;not null"`
}

func InitAPI() (err error) {
	var api API
	var perms []Permission
	err = api.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		return
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
	return nil
}

func (a *API) Migrate() error {
	return apps.Db.AutoMigrate(&a)
}

func getAllApiFromDoc() (res []API, err error) {
	var swagger utils.SwaggerJson
	r := regexp.MustCompile(`"schemes":.*?,`)
	str := r.ReplaceAllString(docs.SwaggerInfo.SwaggerTemplate, "")
	err = json.Unmarshal([]byte(str), &swagger)
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
