package micro_states

type StateMeta struct {
	TtlInSeconds uint64 `json:"ttlInSeconds"`
}

type RedisIdempotent struct {
}

func NewStateMeta(metas ...StateMeta) *StateMeta {
	meta := &StateMeta{}
	for _, m := range metas {
		if m.TtlInSeconds > 0 {
			meta.TtlInSeconds = m.TtlInSeconds
		}
	}
	return meta
}

func (m *StateMeta) SetTtlInSeconds(value uint64) *StateMeta {
	m.TtlInSeconds = value
	return m
}

func (m *StateMeta) ToMap() map[string]string {
	metaMap := make(map[string]string)
	if m.TtlInSeconds > 0 {
		metaMap["ttlInSeconds"] = string(m.TtlInSeconds)
	}
	return metaMap
}
