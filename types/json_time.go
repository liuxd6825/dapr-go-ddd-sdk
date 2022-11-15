package types

import "time"

type JSONTime time.Time

type JSONTimeOption struct {
	time *time.Time
}

var (
	timeJSONFormat = "2006-01-02 15:04:05"
)

func NewJSONTime(value ...*time.Time) *JSONTime {
	var v JSONTime
	if len(value) == 0 {
		v = JSONTime(time.Now())
	} else {
		t := *value[0]
		v = JSONTime(t)
	}
	return &v
}

func AsJSONTime(value *time.Time) *JSONTime {
	if value == nil {
		return nil
	}
	t := *value
	v := JSONTime(t)
	return &v
}

func Now() *JSONTime {
	v := JSONTime(time.Now())
	return &v
}

func SetTimeJSONFormat(format string) {
	timeJSONFormat = format
}

func GetTimeJSONFormat() string {
	return timeJSONFormat
}

func (t *JSONTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeJSONFormat+`"`, string(data), time.Local)
	*t = JSONTime(now)
	return
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte(""), nil
	}
	b := make([]byte, 0, len(timeJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, timeJSONFormat)
	b = append(b, '"')
	return b, nil
}

func (t *JSONTime) String() string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format(timeJSONFormat)
}
