package response

import (
	"gin-bee/apps"
)

type Log struct {
	Id           uint         `json:"id"`
	CreatedAt    apps.FmtTime `json:"createdAt"`
	UserId       *uint        `json:"userId"`
	Method       string       `json:"method"`
	RemoteIP     string       `json:"RemoteIP"`
	Body         string       `json:"Body"`
	Response     string       `json:"response"`
	ResponseTime float64      `json:"responseTime"`
	FullPath     string       `json:"fullPath"`
	Status       string       `json:"status"`
}
type LogResponse struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int   `json:"total"`
	Count    int   `json:"count"`
	Data     []Log `json:"data"`
}
