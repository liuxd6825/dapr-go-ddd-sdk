package goplus

import "time"

var pFalse bool = false
var pTrue bool = true

func PFalse() *bool {
	return &pFalse
}

func PTrue() *bool {
	return &pTrue
}

//////////////
//   Bool
//////////////

func PBool(v bool) *bool {
	return &v
}

func Bool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

//////////////
//   Time
//////////////

func PTime(v time.Time) *time.Time {
	return &v
}

func Time(v *time.Time, deValue time.Time) time.Time {
	if v == nil {
		return deValue
	}
	return *v
}

//////////////
//   String
//////////////

func PString(v string) *string {
	return &v
}

func String(v *string, defVal string) string {
	if v == nil {
		return defVal
	}
	return *v
}

//////////////
//   Int
//////////////

func PInt(v int) *int {
	return &v
}

func Int(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

//////////////
//   Int32
//////////////

func PInt32(v int32) *int32 {
	return &v
}

func Int32(v *int32) int32 {
	if v == nil {
		return 0
	}
	return *v
}

//////////////
//   PInt64
//////////////

func PInt64(v int64) *int64 {
	return &v
}

func Int64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

//////////////
//   PFloat
//////////////

func PFloat(v float32) *float32 {
	return &v
}

func Float(v *float32) float32 {
	if v == nil {
		return 0
	}
	return *v
}

//////////////
//   PFloat64
//////////////

func PFloat64(v float64) *float64 {
	return &v
}

func Float64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
