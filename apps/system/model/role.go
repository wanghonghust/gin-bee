package model

import (
	"gin-bee/apps"
	"gin-bee/zaplog"
)

type Role struct {
	apps.Model
	Name       string       `json:"name" gorm:"type:varchar(64);unique;required;not null"`
	Permission []Permission `json:"permission" gorm:"many2many:role_permissions;"`
	User       []User       `json:"user" gorm:"many2many:user_roles;"`
	Menu       []Menu       `json:"menu" gorm:"many2many:role_menus;"`
}

func InitRolePerM() (err error) {
	var role Role
	err = role.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		return
	}
	zaplog.Logger.Info("数据表role、permission、role_permissions迁移成功")
	return nil
}

func (r *Role) Migrate() error {
	return apps.Db.AutoMigrate(&r)
}
