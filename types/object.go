package types

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gookit/goutil/maputil"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
)

type Object map[string]any

func NewObject() Object {
	return make(map[string]any)
}

func AsObject(val any) (Object, bool) {
	if val == nil {
		return NewObject(), true
	}
	obj, ok := val.(map[string]any)
	if !ok {
		return obj, false
	}
	return Object(obj), true
}

func (o Object) Set(key string, val any) error {
	o[key] = val
	return nil
}

func (o Object) AsMap() map[string]any {
	return map[string]any(o)
}

func (o Object) HasKey(key string) bool {
	return maputil.HasKey(o, key)
}

func (o Object) Get(key string) any {
	return maputil.DeepGet(o, key)
}

func (o Object) GetString(key string) string {
	return convertor.ToString(o[key])
}

func (o Object) GetInt(key string, notNull bool) (*int64, error) {
	var res *int64
	val := o.Get(key)
	if !notNull && (val == nil || val == "") {
		return res, nil
	}
	v, err := convertor.ToInt(o[key])
	if err != nil {
		return nil, errors.New("%s is not int", key)
	}
	return &v, nil
}

func (o Object) GetFloat(key string, notNull bool) (*float64, error) {
	var res *float64
	val := o.Get(key)
	if !notNull && (val == nil || val == "") {
		return res, nil
	}
	v, err := convertor.ToFloat(o[key])
	if err != nil {
		return nil, errors.New("%s is not float64", key)
	}
	return &v, nil

}

func (o Object) GetBool(key string, notNull bool) (*bool, error) {
	var res *bool
	val := o.Get(key)
	if !notNull && (val == nil || val == "") {
		return res, nil
	}
	b, err := convertor.ToBool(o.GetString(key))
	if err != nil {
		return nil, errors.New("%s is not bool", key)
	}
	return &b, nil
}

func (o Object) GetDate(key string) (*times.Date, error) {
	val := o.Get(key)
	return times.AsDate(val)
}

func (o Object) GetTime(key string) (*times.Time, error) {
	val := o.Get(key)
	return times.AsTime(val)
}

func (o Object) GetMap(key string) (map[string]any, bool) {
	val, ok := o[key]
	if ok {
		obj, ok := val.(map[string]any)
		return obj, ok
	}
	return nil, false
}

func (o Object) ToString() string {
	return maputil.ToString(o)
}

func (o Object) ValueOf() any {
	return map[string]any(o)
}

func (o Object) GetData() any {
	return o.Get("data")
}

func (o Object) Keys() []string {
	return maputils.GetKeys(o)
}
