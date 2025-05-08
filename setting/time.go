package setting

import (
	"time"
)

// TimeZone 时区
var timeLocal = time.Local
var timeUTC = time.UTC
var timeZone = timeUTC
var isUTCTime = false

const LocalDateFormat = "2006-01-02"
const LocalTimeFormat = "2006-01-02 15:04:05"
const LocalMsTimeFormat = "2006-01-02 15:04:05.000000"

// SetUTCTimeZone
//
//	@Description: 设置为UTC时区
func SetUTCTimeZone() {
	timeZone = timeUTC
	isUTCTime = true
}

// SetLocalTimeZone
//
//	@Description: 设置为本地时区
func SetLocalTimeZone() {
	timeZone = timeLocal
	isUTCTime = false
}

// GetTimeZone
//
//	@Description: 获取当前时区
//	@return *time.Location
func GetTimeZone() *time.Location {
	return timeZone
}

func GetDateLayout() string {
	return LocalDateFormat
}

func GetTimeLayout() string {
	return LocalTimeFormat
}

func GetMsTimeLayout() string {
	return LocalMsTimeFormat
}

// IsLocalTimeZone
//
//	@Description: 是否为本地时区
//	@return bool
func IsLocalTimeZone() bool {
	return !isUTCTime
}

// IsUTCTimeZone
//
//	@Description: 是UTC时区
//	@return bool
func IsUTCTimeZone() bool {
	return isUTCTime
}
