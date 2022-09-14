package timeutils

import (
	"reflect"
	"time"
)

func Now() time.Time {
	t := time.Now()
	v := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	return v
}

func Equal(t1, t2 interface{}) bool {
	n1 := isNil(t1)
	n2 := isNil(t2)
	if n1 && n2 {
		return true
	} else if n1 && !n2 {
		return false
	} else if !n1 && n2 {
		return false
	}

	switch t1.(type) {
	case time.Time:
		{
			v := t1.(time.Time)
			return equalTime(v, t2)
		}
	case *time.Time:
		{
			v := t1.(*time.Time)
			return equalTime(*v, t2)
		}
	}
	return false
}

func equalTime(t1 time.Time, t2 interface{}) bool {
	switch t2.(type) {
	case time.Time:
		{
			v := t2.(time.Time)
			return t1.Equal(v)
		}
	case *time.Time:
		{
			v := t2.(*time.Time)
			return t1.Equal(*v)
		}
	}
	return false
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}
