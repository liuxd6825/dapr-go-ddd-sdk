package times

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Date time.Time

var (
	dateJSONFormat = "2006-01-02"
	dateLocation   = time.Local
)

type IDate interface {
	Date() time.Time
}

func GetDate(value ...*time.Time) *Date {
	for _, v := range value {
		if v != nil {
			t := *v
			d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, dateLocation)
			res := Date(d)
			return &res
		}
	}
	return nil
}

func NewDate() *Date {
	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, dateLocation)
	res := Date(d)
	return &res
}

func NowDate() *Date {
	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, dateLocation)
	res := Date(d)
	return &res
}

func NewDateWithString(val string) (t *Date, err error) {
	if val == "" {
		return nil, errors.New("time value is empty")
	}
	now, err := time.ParseInLocation(`"`+dateJSONFormat+`"`, val, dateLocation)
	*t = Date(now)
	return t, err
}

func SetDateJSONFormat(format string) {
	dateJSONFormat = format
}

func GetDateJSONFormat() string {
	return dateJSONFormat
}

func (t *Date) GetSchemaType() string {
	return "date"
}

func (t *Date) IsSchemaDateTime() bool {
	return true
}

/*func (t *Date) MarshalJSON() ([]byte, error) {
	if t == nil || t.IsNil() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(dateJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(*t).AppendFormat(b, dateJSONFormat)
	b = append(b, '"')
	return b, nil
}*/

/*
	func (t *Date) UnmarshalJSON(data []byte) (err error) {
		now, err := time.ParseInLocation(`"`+dateJSONFormat+`"`, string(data), time.Local)
		if err != nil {
			return err
		}
		*t = Date(now)
		return err
	}
*/

func (t Date) MarshalBSONValue() (bsontype.Type, []byte, error) {
	tt := time.Time(t)
	return bson.MarshalValue(tt)
}

var shanghaiLocation, _ = time.LoadLocation("Asia/Shanghai")

func (t *Date) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	var tt time.Time
	if err := bson.UnmarshalValue(bt, data, &tt); err != nil {
		return err
	}
	if tt.Location() == time.UTC {
		tt.In(shanghaiLocation)
	}
	*t = Date(tt)
	return nil
}

func (t Date) String() string {
	return time.Time(t).Format(dateJSONFormat)
}

func (t *Date) PTime() *time.Time {
	if t == nil {
		return nil
	}
	v := time.Time(*t)
	return &v
}

func (t *Date) PDate() *time.Time {
	return t.PTime()
}

func (t *Date) Date() time.Time {
	return t.Time()
}

func (t *Date) Time() time.Time {
	if t == nil {
		return time.Time{}
	}
	v := time.Time(*t)
	return v
}

func (t *Date) IsNil() bool {
	return t == nil
}
