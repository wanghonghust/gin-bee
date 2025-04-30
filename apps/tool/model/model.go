package model

import (
	"gin-bee/apps"
	system "gin-bee/apps/system/model"
	"gin-bee/zaplog"
)

type Task struct {
	apps.Model
	Creator      *uint         `json:"creator" gorm:"not null"`
	Uid          string        `json:"uid" gorm:"type:varchar(128);unique;not null"`
	User         system.User   `json:"-" gorm:"foreignKey:Creator;association_foreignkey:ID"`
	Name         string        `json:"name" gorm:"type:varchar(64);unique;not null"`
	RegisterName string        `json:"registerName" gorm:"type:varchar(64);not null"`
	Time         *apps.FmtTime `json:"time"`
	Type         uint          `json:"type" gorm:"type:smallint;not null"`
	TZone        string        `json:"TZone" gorm:"type:varchar(64);"`
	Desc         string        `json:"desc" gorm:"type:varchar(1000)"`
	State        string        `json:"state" gorm:"type:varchar(16);"`
	Result       string        `json:"result" gorm:"type:varchar(256);"`
}

func InitTask() (err error) {
	task := Task{}
	err = task.Migrate()
	if err != nil {
		zaplog.Logger.Error("数据表tasks迁移失败")
		return
	}
	zaplog.Logger.Info("数据表tasks迁移成功")
	return nil
}

func (t *Task) Migrate() error {
	return apps.Db.AutoMigrate(&t)
}
