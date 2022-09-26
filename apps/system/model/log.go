package model

import (
	"gin-bee/apps"
	"gin-bee/zaplog"
)

type Log struct {
	apps.Model
	UserId       *uint   `json:"userId" gorm:"not null"`
	User         User    `json:"-" gorm:"foreignKey:UserId;association_foreignkey:ID"`
	Method       string  `json:"method" gorm:"type:varchar(64);not null"`
	RemoteIP     string  `json:"remoteIP" gorm:"type:varchar(64);not null"`
	Body         string  `json:"body" gorm:"type:text;"`
	Response     string  `json:"response" gorm:"type:mediumblob;"`
	ResponseTime float64 `json:"responseTime" gorm:"not null;"`
	FullPath     string  `json:"fullPath" gorm:"type:varchar(128); not null"`
	Status       int     `json:"status" gorm:"type:smallint;not null"`
}

func InitLog() (err error) {
	log := Log{}
	err = log.Migrate()
	if err != nil {
		zaplog.Logger.Error("数据表logs迁移失败")
		return
	}
	zaplog.Logger.Info("数据表logs迁移成功")
	return nil
}

func (l *Log) Migrate() error {
	return apps.Db.AutoMigrate(&l)
}
