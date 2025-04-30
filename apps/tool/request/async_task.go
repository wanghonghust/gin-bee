package request

import (
	"gin-bee/apps"
)

type AddParam struct {
	Name  string        `json:"name"`
	Time  *apps.FmtTime `json:"time"`
	Type  uint          `json:"type"`
	TZone string        `json:"TZone" `
	Desc  string        `json:"desc"`
}
