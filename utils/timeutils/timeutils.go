package timeutils

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const LocalTimeLayoutLine = "2006-01-02 15:04:05"
const LocalTimeLayoutSlash = "2006/01/02 15:04:05"

const LocalDateLayoutLine = "2006-01-02"
const LocalDateLayoutSlash = "2006/01/02"

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
func PNow() *time.Time {
	t := time.Now()
	return &t
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
		return StrToDateTime(data.(string))
	case *string:
		str := data.(*string)
		if str == nil {
			return defaultValue, nil
		}
		return StrToDateTime(*str)
	case float64:
		return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
	case int64:
		return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
	}
	return defaultValue, nil
}

// 20180313114933

/**/
func StrToDateTime(str string) (time.Time, error) {
	format := LocalTimeLayoutLine
	if len(str) == 14 {
		if res, err := NumStrToDate(str); err == nil {
			return res, nil
		}
	}
	if len(str) <= 10 {
		var err error
		str, err = asDateString(str)
		if err != nil {
			return time.Time{}, nil
		}
		format = LocalDateLayoutLine

	} else {
		if strings.Contains(str, "T") {
			format = time.RFC3339
		} else if strings.Contains(str, "-") {
			format = LocalTimeLayoutLine
		} else if strings.Contains(str, "/") {
			format = LocalTimeLayoutSlash
		} else if strings.Contains(str, "Z") {
			format = time.RFC3339Nano
		}
	}
	res, err := time.Parse(format, str)
	return res, err
}

func asDateString(str string) (string, error) {
	sep := ""
	if strings.Contains(str, "-") {
		sep = "-"
	} else if strings.Contains(str, "/") {
		sep = "/"
	}
	s := strings.Split(str, sep)
	if len(s) != 3 {
		return "", errors.New("error")
	}
	if len(s[1]) == 1 {
		s[1] = "0" + s[1]
	}
	if len(s[2]) == 1 && len(s[0]) == 4 {
		s[2] = "0" + s[2]
	}
	return s[0] + "-" + s[1] + "-" + s[2], nil
}

func NumStrToDate(str string) (t time.Time, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	if _, err := strconv.ParseInt(str, 10, 64); err != nil {
		return time.Time{}, err
	}
	year, err := subToInt(str, 0, 3)
	if err != nil {
		return time.Time{}, err
	}
	month, err := subToInt(str, 4, 5)
	if err != nil {
		return time.Time{}, err
	}
	day, err := subToInt(str, 6, 7)
	if err != nil {
		return time.Time{}, err
	}
	h, err := subToInt(str, 8, 9)
	if err != nil {
		return time.Time{}, err
	}
	m, err := subToInt(str, 10, 11)
	if err != nil {
		return time.Time{}, err
	}
	s, err := subToInt(str, 12, 13)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(year, time.Month(month), day, h, m, s, 0, time.Local), nil
}

func StrToTimePart(str string) (time.Time, error) {
	tNil := time.Time{}
	if len(str) == 8 {
		h, err := subToInt(str, 0, 1)
		if err != nil {
			return tNil, err
		}
		m, err := subToInt(str, 2, 3)
		if err != nil {
			return tNil, err
		}
		s, err := subToInt(str, 4, 5)
		if err != nil {
			return tNil, err
		}
		ns, err := subToInt(str, 6, 7)
		if err != nil {
			return tNil, err
		}

		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, ns*10000000, time.Local), nil
	}
	return tNil, errors.New(str + " is error time.")
}

func subToInt(str string, index int, last int) (int, error) {
	v := str[index : last+1]
	i, err := strconv.ParseInt(v, 10, 64)
	return int(i), err
}

//
// ToDateTime
// @Description:
// @param dateStr
// @param timeStr
// @return *time.Time
//
func ToDateTime(dateStr string, timeStr string) time.Time {
	tstr := strings.ToLower(timeStr)
	d, err := StrToDateTime(dateStr)
	if err != nil {
		return time.Time{}
	}
	t, err := StrToTimePart(tstr)
	if err != nil {
		return time.Time{}
	}
	res := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	return res
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
