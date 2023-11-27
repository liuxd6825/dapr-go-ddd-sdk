package intutils

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
