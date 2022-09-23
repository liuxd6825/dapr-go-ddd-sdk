package timeutils

import (
	"reflect"
	"strings"
	"time"
)

const LocalTimeLayout = "2006-01-02 15:04:05"

//
// Now
// @Description: 获取毫秒值为0的当前时间
// @return time.Time
//
func Now() time.Time {
	t := time.Now()
	v := Time(&t)
	return *v
}

//
// Time
// @Description: 获取毫秒值为0的时间
// @param t
// @return *time.Time
//
func Time(t *time.Time) *time.Time {
	if t == nil {
		return t
	}
	v := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return &v
}

func AnyToTime(data interface{}, defaultValue time.Time) (time.Time, error) {
	switch data.(type) {
	case string:
		return StrToTime(data.(string))
	case *string:
		str := data.(*string)
		if str == nil {
			return defaultValue, nil
		}
		return StrToTime(*str)
	case float64:
		return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
	case int64:
		return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
	}
	return defaultValue, nil
}

func StrToTime(str string) (time.Time, error) {
	format := LocalTimeLayout
	if strings.Contains(str, "T") {
		format = time.RFC3339
	}
	res, err := time.Parse(format, str)
	if err != nil {
		res, err = time.Parse(time.RFC3339Nano, str)
	}
	return res, err
}

//
// Equal
// @Description: 对比时间，精确到秒，忽略毫秒；可以解决mongo时间精度与time.Now不同。
// @param t1  nil/time.Time/*time.Time
// @param t2  nil/time.Time/*time.Time
// @return bool 是否相等
//
func Equal(t1, t2 interface{}) bool {
	switch t1.(type) {
	case time.Time:
		{
			v := t1.(time.Time)
			return equalTime(&v, t2)
		}
	case *time.Time:
		{
			v := t1.(*time.Time)
			return equalTime(v, t2)
		}
	}
	return false
}

func equalTime(t1 *time.Time, t2 interface{}) bool {
	var v2 *time.Time
	switch t2.(type) {
	case time.Time:
		{
			v := t2.(time.Time)
			v2 = &v
		}
	case *time.Time:
		{
			v := t2.(*time.Time)
			if t1 == nil && v == nil {
				return true
			}
			if (t1 == nil && v != nil) || (t1 != nil && v == nil) {
				return false
			}
			v2 = v
		}
	}

	if t1.Year() == v2.Year() && t1.Month() == v2.Month() && t1.Day() == v2.Day() {
		if t1.Hour() == v2.Hour() && t1.Minute() == v2.Minute() && t1.Second() == v2.Second() {
			return true
		}
	}
	return false
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
