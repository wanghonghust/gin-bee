package apps

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type FmtTime time.Time

type Model struct {
	ID        uint     `gorm:"primaryKey" json:"id"`
	CreatedAt FmtTime  `json:"createdAt" binding:"-"`
	UpdatedAt FmtTime  `json:"-" binding:"-"`
	DeletedAt *FmtTime `gorm:"index" json:"-" binding:"-"`
}

func (t *FmtTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(*t)
	return []byte(fmt.Sprintf("\"%v\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func (t *FmtTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tim, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	*t = FmtTime(tim)
	return err
}

func (t *FmtTime) ToString() string {
	tTime := time.Time(*t)
	return tTime.Format("2006-01-02 15:04:05")
}

func (t FmtTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	//判断给定时间是否和默认零时间的时间戳相同
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *FmtTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = FmtTime(value.Local())
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
