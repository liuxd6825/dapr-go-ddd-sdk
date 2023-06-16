package types

import "time"

type JSONDate time.Time

var (
	dateJSONFormat = "2006-01-02"
)

func SetDateJSONFormat(format string) {
	dateJSONFormat = format
}

func GetDateJSONFormat() string {
	return dateJSONFormat
}

func (t *JSONDate) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+dateJSONFormat+`"`, string(data), time.Local)
	*t = JSONDate(now)
	return
}

func (t JSONDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(dateJSONFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, dateJSONFormat)
	b = append(b, '"')
	return b, nil
}

func (t JSONDate) String() string {
	return time.Time(t).Format(dateJSONFormat)
}

func (t *JSONDate) PTime() *time.Time{
	if t==nil{
		return nil
	}
	v := time.Time(*t)
	return &v
}

func (t *JSONDate) Time() time.Time{
	if t==nil{
		return time.Time{}
	}
	v := time.Time(*t)
	return v
}

func (t *JSONDate) IsNil() bool{
	return t==nil
}

