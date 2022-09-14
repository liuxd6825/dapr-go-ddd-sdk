package timeutils

import (
	"reflect"
	"time"
)

//
//  Now
//  @Description: 获取毫秒值为0的当前时间
//  @return time.Time
//
func Now() time.Time {
	t := time.Now()
	v := Time(&t)
	return *v
}

//
//  Time
//  @Description: 获取毫秒值为0的时间
//  @param t
//  @return *time.Time
//
func Time(t *time.Time) *time.Time {
	if t == nil {
		return t
	}
	v := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return &v
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
