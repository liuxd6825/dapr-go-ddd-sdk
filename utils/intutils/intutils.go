package intutils

import "strconv"

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
