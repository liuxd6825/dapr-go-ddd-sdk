package singleutils

import (
	"reflect"
	"sync"
)

type onceItem struct {
	once sync.Once
	obj  any
	fun  any
	key  string
}

var lock sync.RWMutex
var items = make(map[any]*onceItem)

type Options struct {
	GetFuncSkip int
}

type CreateOptions func(opt *Options)

func Create[T any](key string, fun func() T, opts ...CreateOptions) T {
	opt := &Options{GetFuncSkip: 1}
	for _, o := range opts {
		o(opt)
	}
	val, find := getItem(key)
	if find {
		return val.obj.(T)
	}
	val = &onceItem{
		key: key,
		fun: fun,
	}
	val.once.Do(func() {
		if newFun, ok := val.fun.(func() T); ok {
			val.obj = newFun()
			addItem(val)
		}
	})
	return val.obj.(T)
}

func CreateObj[T any](fun func() T, opts ...CreateOptions) T {
	var null T
	key := GetTypeName(null)
	return Create[T](key, fun, opts...)
}

func GetTypeName(val any) string {
	t := reflect.TypeOf(val)
	if t == nil {
		panic("val is nil")
	}
	return t.PkgPath() + t.String()
}

func getItem(key string) (*onceItem, bool) {
	lock.RLock()
	defer func() {
		lock.RUnlock()
	}()
	val, find := items[key]
	return val, find
}

func addItem(item *onceItem) {
	lock.Lock()
	defer func() {
		lock.Unlock()
	}()
	items[item.key] = item
}
