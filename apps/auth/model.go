package auth

import (
	"gin-bee/apps"
	auth "gin-bee/apps/system"
	"gorm.io/gorm"
)

type UserModel struct {
	apps.Model
	Username string    `json:"username" gorm:"type:varchar(64);unique;required;not null"`
	Password string    `json:"password" gorm:"type:varchar(256);required;not null;"`
	Nickname string    `json:"nickname" gorm:"type:varchar(64);unique;"`
	Email    string    `json:"email" gorm:"type:varchar(64)"`
	Avatar   *uint     `json:"avatar" gorm:"default null"`
	File     auth.File `json:"file" gorm:"foreignKey:Avatar;association_foreignkey:ID"`
	State    bool      `json:"state" gorm:"default:true;not null;"` //状态，禁用和启用
}

func (u *UserModel) Create() (tx *gorm.DB) {
	return apps.Db.Create(&u)
}
func (u *UserModel) First(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.First(&u, conds...)
}
func (u *UserModel) Take(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Take(&u, conds...)
}
func (u *UserModel) Last(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Last(&u, conds...)
}
func (u *UserModel) Delete(conds ...interface{}) (tx *gorm.DB) {
	return apps.Db.Delete(&u, conds...)
}
func (u *UserModel) Update(column string, value interface{}) (tx *gorm.DB) {
	return apps.Db.Update(column, value)
}
func (u *UserModel) Updates(values interface{}) (tx *gorm.DB) {
	return apps.Db.Updates(values)
}
func (u *UserModel) Migrate() error {
	return apps.Db.AutoMigrate(&u)
}
func (u *UserModel) Where(query interface{}, args ...interface{}) (tx *gorm.DB) {
	return apps.Db.Where(query, args...)
}
