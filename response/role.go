package response

import "gin-bee/apps/system/model"

type RoleData struct {
	model.Role
	Menu   []TreeMenu `json:"menu"`
	MenuId []uint     `json:"menuId"`
}

type RoleResponse struct {
	Data []RoleData
}
