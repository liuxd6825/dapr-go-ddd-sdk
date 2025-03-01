package times

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"time"
)

type Time time.Time

type JSONTimeOption struct {
	time *time.Time
}

type ITime interface {
	Time() time.Time
}

var (
	timeJSONFormat = "2006-01-02 15:04:05"
)

func NewTime(value ...*time.Time) *Time {
	t := Time(time.Now())
	for _, v := range value {
		if v != nil {
			t := Time(*v)
			return &t
		}
	}
	return &t
}

func NewTimeWithString(val string) (t *Time, err error) {
	if val == "" {
		return nil, errors.New("time value is empty")
	}
	now, err := time.ParseInLocation(`"`+timeJSONFormat+`"`, val, time.Local)
	*t = Time(now)
	return t, err
}

func NowTime() *Time {
	v := Time(time.Now())
	return &v
}

func SetTimeJSONFormat(format string) {
	timeJSONFormat = format
}

func GetTimeJSONFormat() string {
	return timeJSONFormat
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeJSONFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte(""), nil
	}
	b := make([]byte, 0, len(timeJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, timeJSONFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(time.Time(t).UTC())
}

func (t *Time) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	return bson.UnmarshalValue(bt, data, t)
}

func (t *Time) GetSchemaType() string {
	return "datetime"
}

func (t *Time) IsSchemaDateTime() bool {
	return true
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	return time.Time(*t).Format(timeJSONFormat)
}

func (t *Time) PTime() *time.Time {
	if t == nil {
		return nil
	}
	v := time.Time(*t)
	return &v
}

func (t *Time) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	v := time.Time(*t)
	return v
}

func (t *Time) Date() *Date {
	if t == nil {
		return nil
	}
	return NewDate(t.PTime())
}

func (t *Time) IsNil() bool {
	return t == nil
}
