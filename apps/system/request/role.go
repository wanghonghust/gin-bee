package request

type AddRoleParam struct {
	Name       string `json:"name"`
	Permission []uint `json:"permission"`
	Menu       []uint `json:"menu"`
}

type EditRoleParam struct {
	Id         uint   `json:"id"`
	Name       string `json:"name"`
	Permission []uint `json:"permission"`
	Menu       []uint `json:"menu"`
}

type DeleteRoleParam struct {
	Id []uint `json:"id"`
}

type RoleMenuParam struct {
	Id uint `json:"id"`
}
