package maskutils

type UpdateMask struct {
	values  map[string]string
	isEmpty bool
}

func NewUpdateMask(v []string) *UpdateMask {
	mask := &UpdateMask{
		values:  make(map[string]string),
		isEmpty: len(v) == 0,
	}
	for _, key := range v {
		mask.values[key] = key
	}
	return mask
}

func (m *UpdateMask) Contain(key string) bool {
	if m.isEmpty {
		return true
	}
	_, ok := m.values[key]
	return ok
}

func (m *UpdateMask) IsEmpty() bool {
	return m.isEmpty
}
