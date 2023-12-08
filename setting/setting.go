package setting

import "time"

// TimeZone 时区
var timeLocal = time.Local
var timeUTC = time.UTC
var timeZone = timeUTC
var isUTCTime = false

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
