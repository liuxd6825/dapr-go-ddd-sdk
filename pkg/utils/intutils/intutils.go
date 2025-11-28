package intutils

import (
	"fmt"
	"strconv"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
	TB = GB * 1024
)

func StrToInt64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}

func P2Int(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func P2Uint(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}

func P2Int32(p *int32) int32 {
	if p == nil {
		return 0
	}
	return *p
}

func P2Uint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}

func P2Int64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func P2Uint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func P2IntDefault(val *int, def int) int {
	if val == nil {
		return def
	}
	return *val
}

// GetFileSizeTitle 将文件大小（以字节为单位）转换为用户友好的字符串
func GetFileSizeTitle(size int64) string {

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2fT", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.2fG", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2fM", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2fK", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}
