package apps

import (
	"time"
)

type Model struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time  `json:"createdAt" binding:"-"`
	UpdatedAt time.Time  `json:"-" binding:"-"`
	DeletedAt *time.Time `gorm:"index" json:"-" binding:"-"`
}
