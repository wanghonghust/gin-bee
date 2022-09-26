package model

import (
	"gin-bee/apps"
	"gin-bee/zaplog"
	"gorm.io/gorm"
)

type File struct {
	apps.Model
	Path string `json:"path" gorm:"unique;not null"`
}

type Menu struct {
	apps.Model
	Label    string `json:"label" gorm:"required;unique;not null;varchar(64);"`
	Parent   *Menu  `json:"-"`
	ParentId *uint  `json:"parentId"`
	Link     string `json:"link" gorm:"varchar(256)"`
	Icon     string `json:"icon" gorm:"varchar(64)"`
	Role     []Role `json:"role" gorm:"many2many:role_menus;"`
	Local    bool   `json:"local"`
}

func InitFileMenu() (err error) {
	file := File{}
	menu := Menu{}
	err = file.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		return
	}
	zaplog.Logger.Info("数据表file迁移成功")

	err = menu.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
		return
	}
	zaplog.Logger.Info("数据表menu迁移成功")
	return nil
}

func (f *File) Migrate() error {
	return apps.Db.AutoMigrate(&f)
}

func (f *File) Create() (tx *gorm.DB) {
	return apps.Db.Create(&f)
}
func (f *File) First(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.First(&f, conds...)
}
func (f *File) Take(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Take(&f, conds...)
}
func (f *File) Last(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Last(&f, conds...)
}
func (f *File) Delete(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Delete(&f, conds...)
}
func (f *File) Update(column string, value interface{}) (tx *gorm.DB) {
	return apps.Db.Update(column, value)
}
func (f *File) Updates(values interface{}) (tx *gorm.DB) {
	return apps.Db.Updates(values)
}
func (f *File) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return apps.Db.Where(query, args...)
}

func (f *Menu) Migrate() error {
	return apps.Db.AutoMigrate(&f)
}

func (f *Menu) Create() (tx *gorm.DB) {
	return apps.Db.Create(&f)
}
func (f *Menu) First(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.First(&f, conds...)
}
func (f *Menu) Take(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Take(&f, conds...)
}
func (f *Menu) Last(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Last(&f, conds...)
}
func (f *Menu) Delete(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Delete(&f, conds...)
}
func (f *Menu) Update(column string, value interface{}) (tx *gorm.DB) {
	return apps.Db.Update(column, value)
}
func (f *Menu) Updates(values interface{}) (tx *gorm.DB) {
	return apps.Db.Updates(values)
}
func (f *Menu) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return apps.Db.Where(query, args...)
}
