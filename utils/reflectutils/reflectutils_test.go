package reflectutils

import (
	"reflect"
	"testing"
)

type Element struct {
	Id string
}

func Test_ReflectUtils(t *testing.T) {

	t.Run("Reflect", func(t *testing.T) {
		var null *Element
		t1 := reflect.TypeOf(null)
		elem := reflect.New(t1.Elem())
		t.Logf("%v", elem.Interface())
	})

	t.Run("NewStruct", func(t *testing.T) {
		s, err := NewStruct[*Element]()
		if err != nil {
			t.Error(err)
			return
		} else {
			s.Id = "1"
		}
	})

	t.Run("NewSliceItemType1", func(t *testing.T) {
		var list []Element
		t1, err := NewSliceItemType(list)
		if err != nil {
			t.Error(err)
			return
		} else if t1 == nil {
			t.Errorf("itemType is nil")
			return
		}
	})

	t.Run("NewSliceItemType2", func(t *testing.T) {
		var list []Element
		t2, err := NewSliceItemType(&list)
		if err != nil {
			t.Error(err)
			return
		} else if t2 == nil {
			t.Errorf("itemType is nil")
			return
		}
		v2 := reflect.New(t2)
		if v2.Kind() == reflect.Ptr {
			v2 = v2.Elem()
		}
		v2.FieldByName("Id").SetString("1")
		t.Logf("v2:%s", v2.Type().Name())

	})

	t.Run("NewSlice", func(t *testing.T) {
		v3list, err := NewSlice[[]*Element]()
		if err != nil {
			t.Error(err)
			return
		}
		v3list = append(v3list, &Element{})
		t.Logf("list.lenth = %v", len(v3list))
	})

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
