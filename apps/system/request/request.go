package request

import (
	"errors"
	"unicode/utf8"
)

type MenuAddParam struct {
	Label          string  `json:"label"`
	ParentId       uint    `json:"parentId"`
	Link           string  `json:"link"`
	Icon           string  `json:"icon"`
	PermissionSign *string `json:"permissionSign"`
	Local          bool    `json:"local"`
}

type MenuEditParam struct {
	Id             uint    `json:"id"` //必传参数
	Label          string  `json:"label"`
	ParentId       uint    `json:"parentId"`
	Link           string  `json:"link"`
	Icon           string  `json:"icon"`
	PermissionSign *string `json:"permissionSign"`
	Local          bool    `json:"local"`
}

func (m *MenuAddParam) Validator() error {
	if m.Label == "" {
		return errors.New("菜单名不能为空")
	} else if utf8.RuneCountInString(m.Label) > 64 {
		return errors.New("菜单名长度不能大于64")
	}
	return nil
}

func (m *MenuEditParam) Validator() error {
	if utf8.RuneCountInString(m.Label) > 64 {
		return errors.New("菜单名长度不能大于64")
	}
	return nil
}
