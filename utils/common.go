package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type iString interface {
	String() string
}
type wrapper struct {
	time.Time
}
type Time struct {
	wrapper
}
type iError interface {
	Error() string
}

func String(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	case Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		// Empty checks.
		if value == nil {
			return ""
		}
		if f, ok := value.(iString); ok {
			// If the variable implements the String() interface,
			// then use that interface to perform the conversion
			return f.String()
		}
		if f, ok := value.(iError); ok {
			// If the variable implements the Error() interface,
			// then use that interface to perform the conversion
			return f.Error()
		}
		// Reflect checks.
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return String(rv.Elem().Interface())
		}
		// Finally, we use json.Marshal to convert.
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

func AtouiPoint(num string) (res uint) {
	tmp, _ := strconv.Atoi(num)
	utmp := uint(tmp)
	res = utmp
	return
}

func StrTimeFormat(strTime string) string {
	timeLayout := "2006-01-02T15:04:05.99+08:00"
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, strTime, loc)
	fmtLayout := "2006-01-02 15:04:05"
	return theTime.Format(fmtLayout)
}

func StructToMap(obj any) (res map[string]any, err error) {
	marshal, err := json.Marshal(obj)
	fmt.Println(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshal, &res)
	if err != nil {
		return nil, err
	}
	return
}

func AnyToUintPtr(obj any) (res uint, err error) {
	switch val := obj.(type) {
	case uint:
		res = val
	case uint8:
		res = uint(val)
	case uint16:
		res = uint(val)
	case uint32:
		res = uint(val)
	case uint64:
		res = uint(val)
	case int:
		res = uint(val)
	case int8:
		res = uint(val)
	case int16:
		res = uint(val)
	case int32:
		res = uint(val)
	case int64:
		res = uint(val)
	case float32:
		res = uint(val)
	case float64:
		res = uint(val)
	case string:
		tmp, err1 := strconv.Atoi(val)
		if err1 != nil {
			res = 0
			err = err1
		}
		res = uint(tmp)
	default:
		res = 0
		err = errors.New(fmt.Sprintf("Cannot convert type %s to uint", reflect.TypeOf(obj)))
	}
	return
}
