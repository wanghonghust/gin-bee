package model

import "gin-bee/apps"

type Permission struct {
	apps.Model
	Type string `json:"type" gorm:"type:varchar(64);required;not null;uniqueIndex:perm;"`
	Name string `json:"name" gorm:"type:varchar(64);required;not null;uniqueIndex:perm;"`
	Desc string `json:"desc" gorm:"type:varchar(1024);default:''"`
	Role []Role `json:"role" gorm:"many2many:role_permissions;"`
}
