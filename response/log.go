package response

import "time"

type Log struct {
	Id           uint      `json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UserId       *uint     `json:"userId"`
	Method       string    `json:"method"`
	RemoteIP     string    `json:"RemoteIP"`
	Body         string    `json:"Body"`
	Response     string    `json:"response"`
	ResponseTime float64   `json:"responseTime"`
	FullPath     string    `json:"fullPath"`
	Status       string    `json:"status"`
}
type LogResponse struct {
	Data []Log `json:"data"`
}
