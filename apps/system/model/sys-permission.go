package model

import "gin-bee/apps"

type Permission struct {
	apps.Model
	Name string `json:"name" gorm:"type:varchar(64);unique;required;not null"`
}
