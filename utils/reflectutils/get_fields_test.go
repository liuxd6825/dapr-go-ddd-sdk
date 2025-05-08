package reflectutils

import (
	"testing"
	"time"
)

type Super struct {
	F1 string
	F2 int
	F3 time.Time
}

type Object struct {
	Super
	Str string
	Int int64
}

func TestGetFields(t *testing.T) {
	obj := &Object{}
	fields, err := GetFields(obj)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(fields.Names())
	if f, ok := fields.Item("F1"); ok {
		t.Log(f.Name, f.Type)
	}
}
