package setting

import "time"

// TimeZone 时区
var timeZone = time.UTC

// SetUTCTimeZone
//
//	@Description: 设置为UTC时区
func SetUTCTimeZone() {
	timeZone = time.UTC
}

// SetLocalTimeZone
//
//	@Description: 设置为本地时氏
func SetLocalTimeZone() {
	timeZone = time.Local
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
	if timeZone == time.UTC {
		return false
	}
	return true
}

// IsUTCTimeZone
//
//	@Description: 是UTC时区
//	@return bool
func IsUTCTimeZone() bool {
	if timeZone == time.UTC {
		return true
	}
	return false
}
