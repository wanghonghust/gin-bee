package request

import "time"

type AddParam struct {
	Name  string     `json:"name"`
	Time  *time.Time `json:"time"`
	Type  uint       `json:"type"`
	TZone string     `json:"TZone" `
	Desc  string     `json:"desc"`
}
