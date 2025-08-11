package times

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

func GetTime(value ...*time.Time) *Time {
	var t *Time
	for _, v := range value {
		if v != nil {
			t := Time(v.UTC())
			return &t
		}
	}
	return t
}

func NewTime() *Time {
	val := time.Now()
	t := GetTime(&val)
	return t
}

func NewTimeWithString(val string) (t *Time, err error) {
	if val == "" {
		return nil, errors.New("time value is empty")
	}
	now, err := time.ParseInLocation(`"`+timeJSONFormat+`"`, val, time.UTC)
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

/*func (t *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	format := timeJSONFormat
	if strings.Contains(str, "T") && strings.Contains(str, "+") {
		//"2006-01-02 15:04:05"
		format = time.RFC3339Nano
	} else if strings.Contains(str, "T") && strings.Contains(str, "Z") {
		format = time.RFC3339
	}
	newTime, err := time.ParseInLocation(`"`+format+`"`, str, time.Local)
	if err != nil {
		return err
	}
	val := GetTime(&newTime)
	*t = *val
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil || t.IsNil() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(timeJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, timeJSONFormat)
	b = append(b, '"')
	return b, nil
}
*/

// MarshalBSONValue 实现bson自定义序列化
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(t.Time())
}

// UnmarshalBSONValue 实现bson自定义反序列化
func (t *Time) UnmarshalBSONValue(bType bsontype.Type, data []byte) error {
	var tt time.Time
	err := bson.UnmarshalValue(bType, data, &tt)
	if err != nil {
		return err
	}
	*t = Time(tt)
	return nil
}

func (t *Time) Equals(v *Time) bool {
	if t == nil && v == nil {
		return true
	}
	v1 := t.Time()
	v2 := v.Time()
	return v1.Equal(v2)
}
func (t *Time) GetSchemaType() string {
	return "datetime"
}

func (t *Time) IsSchemaDateTime() bool {
	return true
}

func (t Time) String() string {
	return time.Time(t).Format(timeJSONFormat)
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
	return GetDate(t.PTime())
}

func (t *Time) IsNil() bool {
	if t == nil {
		return true
	}
	return t.Time().IsZero()
}
