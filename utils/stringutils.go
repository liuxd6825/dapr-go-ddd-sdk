package utils

import "strconv"

func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}
