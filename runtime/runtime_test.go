package runtime

import (
	"reflect"
	"testing"
)

type Element struct {
	Id string
}

func Test_NewSliceItemType(t *testing.T) {
	var list []Element

	t1 := NewSliceItemType(list)
	if t1 == nil {
		t.Errorf("itemType is nil")
	}
	v1 := reflect.New(t1)
	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}
	v1.FieldByName("Id").SetString("1")
	t.Logf("v1:%s", v1.Type().Name())

	t2 := NewSliceItemType(&list)
	if t2 == nil {
		t.Errorf("itemType is nil")
	}
	v2 := reflect.New(t2)
	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}
	v2.FieldByName("Id").SetString("1")
	t.Logf("v2:%s", v2.Type().Name())
}

func Test_SetSlice(t *testing.T) {
	var list []Element
	sourceList := []string{"id1", "id2"}
	if err := SetSlice(sourceList, list, func(source interface{}, target interface{}) {
		value := source.(string)
		targetValue := target.(*Element)
		targetValue.Id = value
		//target.FieldByName("Id").SetString(value)
	}); err != nil {
		t.Error(err)
	}
	println(list)
}
