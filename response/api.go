package response

import "gin-bee/apps/system/model"

type APIInfos struct {
	Data []model.API `json:"data"`
}

type API struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
}
