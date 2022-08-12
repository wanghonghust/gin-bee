package utils

import (
	"encoding/json"
	"errors"
	"math"
)

type Paginator struct {
	Page    uint          `json:"page"`
	PerPage uint          `json:"perPage"`
	Count   uint          `json:"count"`
	PageNum uint          `json:"pageNum"`
	Data    []interface{} `json:"data"`
	Source  []interface{} `json:"-"`
	Scale   uint          `json:"scale"`
}

func (p *Paginator) PageData(page uint) (res map[string]any, err error) {
	if page == 0 {
		return nil, errors.New("page is not allowed equal to zero")
	}
	start := (page - 1) * p.PerPage
	end := start + p.PerPage
	if end >= uint(len(p.Source)) {
		end = uint(len(p.Source))
	}
	if start >= end {
		p.Data = make([]interface{}, 0)
	} else {
		p.Data = p.Source[start:end]
	}

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
func (p *Paginator) Init(data []interface{}, perPage uint) error {
	if perPage == 0 {
		return errors.New("perPage is not allowed equal to zero")
	}
	p.Count = uint(len(data))
	p.PerPage = perPage
	p.PageNum = uint(math.Ceil(float64(p.Count) / float64(perPage)))
	p.Source = data
	return nil

}
