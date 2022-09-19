package request

import (
	"errors"
	"unicode/utf8"
)

type AddPermissionParam struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (a *AddPermissionParam) Validator() error {
	lenName := utf8.RuneCountInString(a.Name)
	if lenName == 0 || lenName > 64 {
		return errors.New("name不能为空且长度不能大于64")
	}
	lenDesc := utf8.RuneCountInString(a.Desc)
	if lenDesc > 1000 {
		return errors.New("desc长度不能大于1000")
	}
	return nil
}

type EditPermissionParam struct {
	Id   uint   `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (a *EditPermissionParam) Validator() error {
	lenName := utf8.RuneCountInString(a.Name)
	if lenName == 0 || lenName > 64 {
		return errors.New("name不能为空且长度不能大于64")
	}
	lenDesc := utf8.RuneCountInString(a.Desc)
	if lenDesc > 1000 {
		return errors.New("desc长度不能大于1000")
	}
	return nil
}

type DeletePermissionParam struct {
	Id []uint `json:"id"`
}
