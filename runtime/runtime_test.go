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

func Test_MappingSlice(t *testing.T) {
	var targetList []Element
	sourceList := []string{"id1", "id2"}
	if err := MappingSlice(sourceList, targetList, func(i int, source reflect.Value, target reflect.Value) error {
		value := source.Interface().(string)
		targetValue := target.Interface().(*Element)
		targetValue.Id = value
		return nil
	}); err != nil {
		t.Error(err)
	}
	println(targetList)
}
