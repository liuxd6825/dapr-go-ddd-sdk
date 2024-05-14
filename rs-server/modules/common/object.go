package common

import (
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/gookit/goutil/maputil"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"time"
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

func (o Object) GetInt(key string) (int64, error) {
	return convertor.ToInt(o[key])
}

func (o Object) GetFloat(key string) (float64, error) {
	val := o[key]
	return convertor.ToFloat(val)
}

func (o Object) GetBool(key string) (bool, error) {
	return convertor.ToBool(o.GetString(key))
}

func (o Object) GetDateTime(key string) (time.Time, error) {
	return timeutils.AsTime(o.Get(key))
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

func (o Object) Keys() []string {
	return maputils.GetKeys(o)
}
