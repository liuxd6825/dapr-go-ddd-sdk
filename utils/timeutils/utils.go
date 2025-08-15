package timeutils

import "time"
import "google.golang.org/protobuf/types/known/timestamppb"

var timeUtils = NewTimeUtils()

// Now
// @Description: 获取毫秒值为0的当前时间
// @return time.Time
func Now() time.Time {
	return timeUtils.Now()
}
func PNow() *time.Time {
	t := Now()
	return &t
}

// Time
// @Description: 获取毫秒值为0的时间
// @param t
// @return *time.Time
func Time(t *time.Time) *time.Time {
	if t == nil {
		return t
	}
	v := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return &v
}

func AnyToTime(data any, defaultValue time.Time) (res time.Time, err error) {
	return timeUtils.AnyToTime(data, defaultValue)
}

func AsTime(data any) (time.Time, error) {
	return timeUtils.AsTime(data)
}

func AsJsonTime(data any) (time.Time, error) {
	return AnyToTime(data, time.Time{})
}

func StrToDateTime(str string) (res time.Time, err error) {
	return timeUtils.StrToDateTime(str)
}

func FormatStr(fmt, str string) (res string, err error) {
	return timeUtils.FormatStr(fmt, str)
}
func NumStrToDate(str string) (t time.Time, err error) {
	return timeUtils.NumStrToDate(str)
}

func StrToTimePart(v string) (time.Time, error) {
	return timeUtils.StrToTimePart(v)
}

func ToDateTime(dateStr string, timeStr string) (time.Time, error) {
	return timeUtils.ToDateTime(dateStr, timeStr)
}

func Equal(t1, t2 any) bool {
	return timeUtils.Equal(t1, t2)
}

func ToTimeString(date *time.Time) string {
	return timeUtils.ToTimeString(date)
}

func ToDateString(date *time.Time) string {
	return timeUtils.ToDateString(date)
}

func AsTimestamp(t *time.Time) *timestamppb.Timestamp {
	return timeUtils.AsTimestamp(t)
}

func ToPTime(t time.Time) *time.Time {
	return &t
}

func NewDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}
