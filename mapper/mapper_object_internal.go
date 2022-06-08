package mapper

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"reflect"
	"strings"
	"time"
)

// registerValue register Value to init Map
func (dm *mapperObject) registerValue(objValue reflect.Value) error {
	regValue := objValue
	if objValue == dm.ZeroValue {
		return errors.New("obj value does not exist")
	}

	valueType := regValue.Type()
	if valueType.Kind() == reflect.Ptr {
		regValue = regValue.Elem()
	}

	typeName := valueType.String()
	if valueType.Kind() == reflect.Struct {
		for i := 0; i < regValue.NumField(); i++ {
			fieldName := dm.getFieldName(regValue, i)
			if fieldName == IgnoreTagValue {
				continue
			}
			mapFieldName := typeName + nameConnector + fieldName
			if valueType.Field(i).Type.Kind() == reflect.Struct {
				dm.registerBaseValue(typeName, regValue.Field(i), valueType.Field(i).Type)
			} else {
				realFieldName := valueType.Field(i).Name
				dm.fieldNameMap.Store(mapFieldName, realFieldName)
			}
		}
	}

	// store register flag
	dm.registerMap.Store(typeName, nil)
	return nil
}

func (dm *mapperObject) registerBaseValue(typeName string, regValue reflect.Value, regType reflect.Type) {
	for i := 0; i < regValue.NumField(); i++ {
		fieldName := dm.getFieldName(regValue, i)
		if fieldName == IgnoreTagValue {
			continue
		}
		mapFieldName := typeName + nameConnector + fieldName
		if regType.Field(i).Type.Kind() == reflect.Struct {
			dm.registerBaseValue(typeName, regValue.Field(i), regType.Field(i).Type)
		} else {
			realFieldName := regType.Field(i).Name
			dm.fieldNameMap.Store(mapFieldName, realFieldName)
		}
	}
}

// GetFieldName get fieldName with ElemValue and index
// if config tag string, return tag value
func (dm *mapperObject) getFieldName(objElem reflect.Value, index int) string {
	fieldName := ""
	field := objElem.Type().Field(index)
	tag := dm.getStructTag(field)

	// keeps the behavior in old version
	if tag == IgnoreTagValue && !dm.IsEnableFieldIgnoreTag() {
		tag = ""
	}

	if tag != "" {
		fieldName = tag
	} else {
		fieldName = field.Name
	}
	return fieldName
}

// UseWrapper register a type wrapper
func (dm *mapperObject) useWrapper(w TypeWrapper) {
	if len(dm.typeWrappers) > 0 {
		dm.typeWrappers[len(dm.typeWrappers)-1].SetNext(w)
	}
	dm.typeWrappers = append(dm.typeWrappers, w)
}

func (dm *mapperObject) elemMapper(fromElem, toElem reflect.Value) error {
	// check register flag
	// if not register, register it
	if !dm.checkIsRegister(fromElem) {
		if err := dm.registerValue(fromElem); err != nil {
			return err
		}
	}
	if !dm.checkIsRegister(toElem) {
		if err := dm.registerValue(toElem); err != nil {
			return err
		}
	}
	if toElem.Type().Kind() == reflect.Map {
		return dm.elemToMap(fromElem, toElem)
	}

	return dm.elemToStruct(fromElem, toElem)
}

func (dm *mapperObject) elemToStruct(fromElem, toElem reflect.Value) (resErr error) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				msg := fmt.Sprintf("mapper %s to %s error:  %s ", fromElem.Type().Name(), toElem.Type().Name(), err.Error())
				resErr = errors.New(msg)
			}
		}
	}()
	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := dm.getFieldName(fromElem, i)

		if fieldName == IgnoreTagValue {
			continue
		}

		// check field is exists
		realFieldName, exists := dm.CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}

		toFieldInfo := toElem.FieldByName(realFieldName)
		// check field is same type
		if dm.enabledTypeChecking {
			if fromFieldInfo.Kind() != toFieldInfo.Kind() {
				continue
			}
		}

		if dm.enabledMapperStructField &&
			toFieldInfo.Kind() == reflect.Struct && fromFieldInfo.Kind() == reflect.Struct &&
			toFieldInfo.Type() != fromFieldInfo.Type() &&
			!dm.checkIsTypeWrapper(toFieldInfo) && !dm.checkIsTypeWrapper(fromFieldInfo) {
			x := reflect.New(toFieldInfo.Type()).Elem()
			err := dm.elemMapper(fromFieldInfo, x)
			if err != nil {
				fmt.Println("auto mapper failed", fromFieldInfo, "=>", toFieldInfo, "error", err.Error())
			} else {
				toFieldInfo.Set(x)
			}
		} else {
			isSet := false
			if dm.enabledAutoTypeConvert {
				for _, typeWrapper := range dm.typeWrappers {
					set, err := typeWrapper.SetValue(fromFieldInfo, toFieldInfo)
					if err != nil {
						return err
					} else if set {
						isSet = true
						break
					}
				}
			}
			if !isSet {
				toFieldInfo.Set(fromFieldInfo)
			}
		}
	}
	return nil
}

func (dm *mapperObject) elemToMap(fromElem, toElem reflect.Value) error {
	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := dm.getFieldName(fromElem, i)
		if fieldName == IgnoreTagValue {
			continue
		}
		toElem.SetMapIndex(reflect.ValueOf(fieldName), fromFieldInfo)
	}
	return nil
}

func (dm *mapperObject) setFieldValue(fieldValue reflect.Value, fieldKind reflect.Kind, value interface{}) error {
	switch fieldKind {
	case reflect.Bool:
		if value == nil {
			fieldValue.SetBool(false)
		} else if v, ok := value.(bool); ok {
			fieldValue.SetBool(v)
		} else {
			v, _ := Convert(ToString(value)).Bool()
			fieldValue.SetBool(v)
		}

	case reflect.String:
		if value == nil {
			fieldValue.SetString("")
		} else {
			fieldValue.SetString(ToString(value))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			fieldValue.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetInt(int64(val.Uint()))
			default:
				v, _ := Convert(ToString(value)).Int64()
				fieldValue.SetInt(v)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			fieldValue.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(val.Uint())
			default:
				v, _ := Convert(ToString(value)).Uint64()
				fieldValue.SetUint(v)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			fieldValue.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				fieldValue.SetFloat(val.Float())
			default:
				v, _ := Convert(ToString(value)).Float64()
				fieldValue.SetFloat(v)
			}
		}
	case reflect.Struct:
		if value == nil {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
		} else if dm.DefaultTimeWrapper.IsType(fieldValue) {
			var timeString string
			if fieldValue.Type() == dm.timeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(value))
			}
			if fieldValue.Type() == dm.jsonTimeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(types.JSONTime(value.(time.Time))))
			}
			switch d := value.(type) {
			case []byte:
				timeString = string(d)
			case string:
				timeString = d
			case int64:
				if dm.enabledAutoTypeConvert {
					// try to transform Unix time to local Time
					t, err := UnixToTimeLocation(value.(int64), time.UTC.String())
					if err != nil {
						return err
					}
					fieldValue.Set(reflect.ValueOf(t))
				}
			}
			if timeString != "" {
				if len(timeString) >= 19 {
					// 满足yyyy-MM-dd HH:mm:ss格式
					timeString = timeString[:19]
					t, err := time.ParseInLocation(formatDateTime, timeString, time.UTC)
					if err == nil {
						t = t.In(time.UTC)
						fieldValue.Set(reflect.ValueOf(t))
					}
				} else if len(timeString) >= 10 {
					// 满足yyyy-MM-dd格式
					timeString = timeString[:10]
					t, err := time.ParseInLocation(formatDate, timeString, time.UTC)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(t))
					}
				}
			}
		}
	default:
		if reflect.ValueOf(value).Type() == fieldValue.Type() {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func (dm *mapperObject) getStructTag(field reflect.StructField) string {
	tagValue := ""
	// 1.check mapperTagKey
	if dm.enabledMapperTag {
		tagValue = field.Tag.Get(mapperTagKey)
		if tagValue != "" {
			return tagValue
		}
	}

	// 2.check jsonTagKey
	if dm.enabledJsonTag {
		tagValue = field.Tag.Get(jsonTagKey)
		if tagValue != "" {
			// support more tag property, as json tag omitempty 2018-07-13
			return strings.Split(tagValue, ",")[0]
		}
	}

	return tagValue
}

func (dm *mapperObject) checkIsRegister(objElem reflect.Value) bool {
	typeName := objElem.Type().String()
	_, isOk := dm.registerMap.Load(typeName)
	return isOk
}

// convertToSlice convert slice interface{} to []interface{}
func (dm *mapperObject) convertToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Ptr {
		if v.Elem().Kind() != reflect.Slice {
			panic("fromSlice arr is not a pointer to a slice")
		}
		v = v.Elem()
	} else {
		if v.Kind() != reflect.Slice {
			panic("fromSlice arr is not a slice")
		}
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

// checkIsTypeWrapper check value is in type wrappers
func (dm *mapperObject) checkIsTypeWrapper(value reflect.Value) bool {
	for _, w := range dm.typeWrappers {
		if w.IsType(value) {
			return true
		}
	}
	return false
}
