package model

import (
	"database/sql/driver"
	"encoding/json"
	"gin-bee/apps"
	"gin-bee/zaplog"
	"gorm.io/gorm"
)

type Limiter struct {
	On    bool `json:"on"`
	Limit uint `json:"limit"`
}

// Value 存储数据的时候转换为字符串
func (l Limiter) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan 读取数据的时候转换为json
func (l *Limiter) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &l)
}

type User struct {
	apps.Model
	Username    string  `json:"username" gorm:"type:varchar(64);unique;required;not null"`
	Password    string  `json:"password" gorm:"type:varchar(256);required;not null;"`
	Nickname    string  `json:"nickname" gorm:"type:varchar(64);"`
	Email       string  `json:"email" gorm:"type:varchar(64)"`
	Avatar      *uint   `json:"avatar" gorm:"default null"`
	File        File    `json:"-" gorm:"foreignKey:Avatar;association_foreignkey:ID"`
	State       bool    `json:"state" gorm:"default:true;not null;"` //状态，禁用和启用
	IsSuperUser bool    `json:"isSuperUser" gorm:"default:false;not null;"`
	Role        []Role  `json:"role" gorm:"many2many:user_roles;" binding:"-"`
	Limiter     Limiter `json:"limiter" gorm:"type:json;"`
}

func InitUser() (err error) {
	defer func() {
		if err == nil {
			zaplog.Logger.Info("数据表user迁移成功")
		} else {
			zaplog.Logger.Error("数据表user迁移失败")
		}

	}()
	user := User{}
	err = user.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		return
	}
	return nil
}

func (u *User) Create() (tx *gorm.DB) {
	return apps.Db.Create(&u)
}
func (u *User) First(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.First(&u, conds...)
}
func (u *User) Take(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Take(&u, conds...)
}
func (u *User) Last(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Last(&u, conds...)
}
func (u *User) Delete(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Delete(&u, conds...)
}
func (u *User) Update(column string, value interface{}) (tx *gorm.DB) {
	return apps.Db.Update(column, value)
}
func (u *User) Updates(values interface{}) (tx *gorm.DB) {
	return apps.Db.Updates(values)
}
func (u *User) Migrate() error {
	return apps.Db.AutoMigrate(&u)
}
func (u *User) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return apps.Db.Where(query, args...)
}
