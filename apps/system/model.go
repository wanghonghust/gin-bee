package system

import (
	"gin-bee/apps"
	"gin-bee/zaplog"
	"gorm.io/gorm"
)

type File struct {
	apps.Model
	Path string `json:"path" gorm:"unique;not null"`
}

func init() {
	file := File{}
	err := file.Migrate()
	if err != nil {
		zaplog.Logger.Error(err.Error())
	}
	zaplog.Logger.Info("数据表file迁移成功")
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
