package response

import (
	"time"
)

type Task struct {
	Creator      *uint      `json:"creator" `
	Uid          string     `json:"uid"`
	Name         string     `json:"name"`
	RegisterName string     `json:"registerName"`
	Time         *time.Time `json:"time"`
	Type         uint       `json:"type"`
	TZone        string     `json:"TZone"`
	Desc         string     `json:"desc"`
	State        string     `json:"state"`
	Result       string     `json:"result"`
}
type TaskResponse struct {
	Data []Task `json:"data"`
}
