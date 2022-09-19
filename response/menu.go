package response

type TreeMenu struct {
	ID       uint       `json:"id"`
	Label    string     `json:"label"`
	ParentId *uint      `json:"parentId"`
	Link     string     `json:"link"`
	Icon     string     `json:"icon"`
	Children []TreeMenu `json:"children"`
	CreateAt string     `json:"createAt"`
	Local    bool       `json:"local"`
}

type MenuResponse struct {
	Menus []TreeMenu
}
