package vis_network

import "time"

// Heads 流水时间范围
type TimeHeader struct {
	MinTime *time.Time `json:"minTime"`
	MaxTime *time.Time `json:"maxTime"`
	Days    int        `json:"days"`
	Years   int        `json:"years"`
	MinYear int        `json:"minYear"`
	MaxYear int        `json:"maxYear"`
}

func NewTimeHeader(minTime time.Time, maxTime time.Time) *TimeHeader {
	heads := &TimeHeader{}
	heads.MinTime = &minTime
	heads.MaxTime = &maxTime
	days := maxTime.Sub(minTime).Hours() / 24
	heads.Days = int(days)
	heads.Years = maxTime.Year() - minTime.Year()
	heads.MinYear = minTime.Year()
	heads.MaxYear = maxTime.Year()
	return heads
}
