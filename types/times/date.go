package times

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"time"
)

type Date time.Time

var (
	dateJSONFormat = "2006-01-02"
)

func NewDate(value ...*time.Time) *Date {
	var res Date
	if len(value) == 0 {
		res = Date(time.Now())
	} else {
		for _, v := range value {
			if v != nil {
				t := *v
				d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
				res = Date(d)
				break
			}
		}
	}
	return &res
}

func SetDateJSONFormat(format string) {
	dateJSONFormat = format
}

func GetDateJSONFormat() string {
	return dateJSONFormat
}

func (t *Date) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dateJSONFormat+`"`, string(data), time.Local)
	*t = Date(now)
	return err
}

func (t Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dateJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, dateJSONFormat)
	b = append(b, '"')
	return b, nil
}

func (t Date) MarshalBSONValue() (bsontype.Type, []byte, error) {
	val := t.Time().Add(8 * time.Hour)
	return bson.MarshalValue(val)
}

func (t *Date) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	return bson.UnmarshalValue(bt, data, t)
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
