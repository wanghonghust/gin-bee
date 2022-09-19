package utils

import (
	"encoding/json"
	"errors"
	"math"
)

type Paginator struct {
	Page    int              `json:"page"`
	PerPage int              `json:"perPage"`
	Count   uint             `json:"count"`
	PageNum uint             `json:"pageNum"`
	Data    []map[string]any `json:"data"`
	Source  []map[string]any `json:"-"`
	Scale   uint             `json:"scale"`
}

func (p *Paginator) PageData(page int) (res map[string]any, err error) {
	// Paginator 只是做数据格式转换，实际分页由gorm查询完成
	//start := 0
	//end := len(p.Source)
	//if page < 0 {
	//	return nil, errors.New("page is not allowed less than zero")
	//} else if page > 0 {
	//	start = (page - 1) * p.PerPage
	//	end = start + p.PerPage
	//	if end >= len(p.Source) {
	//		end = len(p.Source)
	//	}
	//}
	//if start >= end {
	//	p.Data = make([]interface{}, 0)
	//} else {
	//	p.Data = p.Source[start:end]
	//}
	p.Data = p.Source
	p.Scale = uint(len(p.Data))

	marshal, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		return nil, err
	}
	return
}
func (p *Paginator) Init(data []map[string]any, perPage int, count uint) error {
	if perPage < 0 {
		return errors.New("perPage is not allowed less than zero")
	}
	p.Count = count
	p.PerPage = perPage
	p.PageNum = uint(math.Ceil(float64(p.Count) / float64(perPage)))
	p.Source = data
	return nil

}
