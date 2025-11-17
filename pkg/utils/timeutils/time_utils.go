package timeutils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	times2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types/times"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const LocalDateFormatLine = "2006-01-02"
const LocalTimeFormatLine = "2006-01-02 15:04:05"
const LocalMsTimeFormatLine = "2006-01-02 15:04:05.000000"

const LocalTimeFormatSlash = "2006/01/02 15:04:05"
const LocalDateLayoutSlash = "2006/01/02"

var (
	splitsTime = []string{":", ",", " ", "."}
	datesTime  = []string{"-", ",", " ", "."}
)

type TimeUtils struct {
	LocalDateFormatLine   string
	LocalTimeFormatLine   string
	LocalMsTimeFormatLine string
	LocalTimeFormatSlash  string
	LocalDateLayoutSlash  string
	TimeLocal             *time.Location
}

func NewTimeUtils() *TimeUtils {
	local := getShanghaiLocation()
	return &TimeUtils{
		LocalDateFormatLine:   LocalDateFormatLine,
		LocalTimeFormatLine:   LocalTimeFormatLine,
		LocalMsTimeFormatLine: LocalMsTimeFormatLine,
		LocalTimeFormatSlash:  LocalTimeFormatSlash,
		LocalDateLayoutSlash:  LocalDateLayoutSlash,
		TimeLocal:             local,
	}
}

var shanghaiLocation *time.Location
var shanghaiLocationOnce sync.Once

func getShanghaiLocation() *time.Location {
	shanghaiLocationOnce.Do(func() {
		local, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			panic(err)
		}
		shanghaiLocation = local
	})
	return shanghaiLocation
}

func (u *TimeUtils) Now() time.Time {
	now := time.Now()
	return now.In(u.TimeLocal)
}

func (u *TimeUtils) AnyToTime(data any, defVal time.Time) (res time.Time, err error) {
	if data == nil {
		return defVal, nil
	}
	switch data.(type) {
	case times2.Time:
		val := data.(times2.Time)
		res = val.Time()
	case times2.Date:
		val := data.(times2.Date)
		res = val.Time()
	case *times2.Time:
		val := data.(*times2.Time)
		res = val.Time()
	case *times2.Date:
		val := data.(*times2.Date)
		res = val.Time()
	case time.Time:
		res = data.(time.Time)
	case *time.Time:
		if data == nil {
			res = defVal
		} else {
			var timeVal = data.(*time.Time)
			res = *timeVal
		}
	case string:
		res, err = StrToDateTime(data.(string))
	case *string:
		str := data.(*string)
		if str == nil {
			return defVal, nil
		}
		res, err = StrToDateTime(*str)
	case float64:
		res, err = time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
	case int64:
		res, err = time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
	default:
		res, err = time.Time{}, errors.New("Invalid data type")
	}
	if err == nil {
		val := times2.GetTime(&res)
		return val.Time(), nil
	}
	return defVal, nil
}

func (u *TimeUtils) AsTime(data any) (time.Time, error) {
	return u.AnyToTime(data, time.Time{})
}

// HasHMS
// @Description: 是否有时分秒部分
// @receiver u
// @param val
// @return bool
func HasHMS(val time.Time) bool {
	if val.IsZero() {
		return false
	}
	if val.Hour() == 0 && val.Minute() == 0 && val.Second() == 0 {
		return false
	}
	return true
}

// 20180313114933
// 20221001 11:09:22

// StrToDateTime
//
//	@Description:
//	@param str
//	@return time.Time
//	@return error
func (u *TimeUtils) StrToDateTime(str string) (res time.Time, err error) {
	format := u.LocalTimeFormatLine

	// 按空格对字符串进行日期与时间的分割
	dayVal, timeVar := split(str)

	// 对日期部分进行格式化
	dayVal, err = fmtDateStr(dayVal)
	if err != nil {
		return time.Time{}, err
	}

	timeIsNil := false
	if len(timeVar) == 0 {
		timeIsNil = true
	}
	str = dayVal
	if !timeIsNil {
		// 对时间部分进行格式化
		timeVar, err = fmtTimeStr(timeVar)
		if err != nil {
			return time.Time{}, err
		}
		str = dayVal + " " + timeVar
	}

	// 合并为完整时间

	if strings.Contains(str, "-") {
		if timeIsNil {
			format = u.LocalDateFormatLine
		} else {
			format = u.LocalTimeFormatLine
		}
	} else if strings.Contains(str, "T") {
		format = time.RFC3339
	} else if strings.Contains(str, "/") {
		format = LocalTimeFormatSlash
	} else if strings.Contains(str, "Z") {
		format = time.RFC3339Nano
	}

	res, err = time.ParseInLocation(format, str, u.TimeLocal)
	return res, err
}

func (u *TimeUtils) FormatStr(fmt, str string) (res string, err error) {
	tVal, err := StrToDateTime(str)
	if err != nil {
		return "", err
	}
	return tVal.Format(fmt), nil
}

func (u *TimeUtils) NumStrToDate(str string) (t time.Time, err error) {
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

func (u *TimeUtils) StrToTimePart(v string) (time.Time, error) {
	tNil := time.Time{}
	rs := []string{"-", ".", ":", " ", ",", "_", "%"}
	str := v
	for _, r := range rs {
		str = strings.Replace(str, r, "", -1)
	}

	var h, m, s, ns int
	var err error
	count := len(str)
	if count >= 6 {
		h, err = subToInt(str, 0, 1)
		if err != nil {
			return tNil, err
		}
		m, err = subToInt(str, 2, 3)
		if err != nil {
			return tNil, err
		}
		s, err = subToInt(str, 4, 5)
		if err != nil {
			return tNil, err
		}
		ns, err = subToInt(str, 6, count-6)
		if err != nil {
			return tNil, err
		}
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), h, m, s, ns*10000000, time.Local), nil
	}
	return tNil, errors.New(str + " is error time.")
}

// ToDateTime
// @Description:
// @param dateStr
// @param timeStr
// @return *time.Time
func (u *TimeUtils) ToDateTime(dateStr string, timeStr string) (time.Time, error) {
	tstr := strings.ToLower(timeStr)
	d, err := StrToDateTime(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	t, err := StrToTimePart(tstr)
	if err != nil {
		return time.Time{}, err
	}
	res := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local)
	return res, nil
}

// Equal
// @Description: 对比时间，精确到秒，忽略毫秒；可以解决mongo时间精度与time.Now不同。
// @param t1  nil/time.Time/*time.Time
// @param t2  nil/time.Time/*time.Time
// @return bool 是否相等
func (u *TimeUtils) Equal(t1, t2 any) bool {
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

func (u *TimeUtils) ToTimeString(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format(u.LocalTimeFormatLine)
}

func (u *TimeUtils) ToDateString(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format(u.LocalDateFormatLine)
}

func (u *TimeUtils) AsTimestamp(t *time.Time) *timestamppb.Timestamp {
	timestamp := &timestamppb.Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
	return timestamp
}

func (u *TimeUtils) ToPTime(t time.Time) *time.Time {
	return &t
}
func (u *TimeUtils) isNil(i any) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// split
//
//	@Description: 将日期字符串分割成日期与时间两部分
//	@param val
//	@return string
//	@return string
func split(val string) (string, string) {
	val = strings.ReplaceAll(val, "  ", " ")
	val = strings.ReplaceAll(val, "T", " ")
	val = strings.ReplaceAll(val, "Z", " ")
	list := strings.Split(val, " ")
	count := len(list)
	if count == 1 {
		return list[0], ""
	}
	if count >= 2 {
		return list[0], list[1]
	}
	return "", ""
}

// fmtTimeStr
//
//	@Description: 将日期字符串格式化为 HH:mm:ss 格式, 支持格式有：121019 |  12:10:09 | 12.10.10 | 12,10,10 | 12:10:09 | 12 10 19 | 12时10分19秒
//	@Description:
//	@param str
//	@return string
//	@return error
func fmtTimeStr(str string) (string, error) {
	str = strings.ReplaceAll(str, " ", "")

	// 处理不规则时间部分 如： 13.32,27; 13.3227; 1332.27
	str = strings.ReplaceAll(str, ".", ":")
	str = strings.ReplaceAll(str, ",", ":")
	str = strings.ReplaceAll(str, "。", ":")
	str = strings.ReplaceAll(str, "，", ":")
	if strings.Contains(str, ":") {
		mList := strings.Split(str, ":")
		var rList []string
		for _, m := range mList {
			mLen := len(m)
			if mLen == 4 {
				rList = append(rList, m[0:2])
				rList = append(rList, m[2:4])
			} else {
				rList = append(rList, m)
			}
		}
		str = strings.Join(rList, ":")
	}

	if len(str) == 4 {
		ok := true
		for _, sp := range splitsTime {
			if strings.Contains(str, sp) {
				ok = false
				break
			}
		}
		if ok {
			y := str[0:2]
			m := str[2:4]
			return fmt.Sprintf("%s:%s:00", y, m), nil
		}
	} else if len(str) == 6 {
		ok := true
		for _, sp := range splitsTime {
			if strings.Contains(str, sp) {
				ok = false
				break
			}
		}
		if ok {
			y := str[0:2]
			m := str[2:4]
			d := str[4:6]
			return fmt.Sprintf("%s:%s:%s", y, m, d), nil
		}
	}

	str = strings.ReplaceAll(str, "时", ":")
	str = strings.ReplaceAll(str, "分", ":")
	str = strings.ReplaceAll(str, "秒", ":")
	sep := ":"

	for _, sp := range splitsTime {
		if strings.Contains(str, sp) {
			sep = sp
			break
		}
	}
	s := strings.Split(str, sep)
	if len(s) < 3 {
		return "", errors.New(fmt.Sprintf(`%s 时间格式不正确，HH:mm:ss`, str))
	}
	if len(s[1]) == 1 {
		s[1] = "0" + s[1]
	}
	if len(s[2]) == 1 {
		s[2] = "0" + s[2]
	}
	return s[0] + ":" + s[1] + ":" + s[2], nil
}

// fmtDateStr
//
//	@Description: 将日期字符串格式化为 yyyy-MM-dd格式, 支持格式有：2000年10月10日 | 20001010 | 2016-10-9 | 2014/10/10 | 2018.10.3
//	@Description:
//	@param str
//	@return string
//	@return error
func fmtDateStr(str string) (string, error) {
	if len(str) == 8 {
		ok := true
		for _, sp := range datesTime {
			if strings.Contains(str, sp) {
				ok = false
				break
			}
		}
		if ok {
			y := str[0:4]
			m := str[4:6]
			d := str[6:8]
			return fmt.Sprintf("%s-%s-%s", y, m, d), nil
		}
	}
	str = strings.ReplaceAll(str, "年", "-")
	str = strings.ReplaceAll(str, "月", "-")
	str = strings.ReplaceAll(str, "日", "-")

	sep := ""
	for _, sp := range datesTime {
		if strings.Contains(str, sp) {
			sep = sp
			break
		}
	}

	s := strings.Split(str, sep)
	if len(s) < 3 {
		return "", errors.New(fmt.Sprintf(`%s 时间格式不正确`, str))
	}
	if len(s[1]) == 1 {
		s[1] = "0" + s[1]
	}
	if len(s[2]) == 1 && len(s[0]) == 4 {
		s[2] = "0" + s[2]
	}
	return s[0] + "-" + s[1] + "-" + s[2], nil
}

func subToInt(str string, index int, last int) (int, error) {
	if last == 0 {
		return 0, nil
	}
	v := str[index : last+1]
	if v == "" {
		return 0, nil
	}
	i, err := strconv.ParseInt(v, 10, 64)
	return int(i), err
}

func equalTime(t1 *time.Time, t2 any) bool {
	var v2 *time.Time
	switch t2.(type) {
	case time.Time:
		v := t2.(time.Time)
		v2 = &v
	case *time.Time:
		v := t2.(*time.Time)
		if t1 == nil && v == nil {
			return true
		}
		if (t1 == nil && v != nil) || (t1 != nil && v == nil) {
			return false
		}
		v2 = v
	}

	if t1.Year() == v2.Year() && t1.Month() == v2.Month() && t1.Day() == v2.Day() {
		if t1.Hour() == v2.Hour() && t1.Minute() == v2.Minute() && t1.Second() == v2.Second() {
			return true
		}
	}
	return false
}
