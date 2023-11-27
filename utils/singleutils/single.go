package singleutils

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"reflect"
	"sync"
)

type item[T any] struct {
	once     sync.Once
	newFun   func() T
	instance T
}

var _objects = cmap.New()

func newSingleton[T any](key string, newFun func() T) *item[T] {
	return &item[T]{
		newFun: newFun,
	}
}

func Get[T any](new ...func() T) (T, error) {
	var null T
	key, err := getKey[T]()
	if err != nil {
		return null, err
	}
	v, ok := _objects.Get(key)
	if !ok {
		if new != nil && len(new) > 0 {
			Set[T](new[0])
			v, ok = _objects.Get(key)
		}
	}
	if !ok {
		return null, fmt.Errorf("singleutils.Get[T]() \"%s\" key does not exist", key)
	}
	s, _ := v.(*item[T])
	s.once.Do(func() {
		s.instance = s.newFun()
	})
	return s.instance, nil
}

func GetSet[T any](new func() T) T {
	v, err := Get[T](new)
	if err != nil {
		panic(err)
	}
	return v
}

func GetObject[T any]() T {
	v, err := Get[T]()
	if err != nil {
		panic(err)
	}
	return v
}

func Set[T any](new func() T) error {
	key, err := getKey[T]()
	if err != nil {
		return err
	}
	_, ok := _objects.Get(key)
	if ok {
		return fmt.Errorf("singleutils.Set[T](new) error: \"%s\" key does exist", key)
	}
	s := newSingleton[T](key, new)
	_objects.Set(key, s)
	return nil
}

func getKey[T any]() (string, error) {
	var null T
	t := reflect.TypeOf(null)
	if t == nil {
		return "", fmt.Errorf("getKey[T any]() error: t is interface")
	}
	elem := t.Elem()
	key := fmt.Sprintf("%s.%s", elem.PkgPath(), elem.Name())
	if t.Kind() == reflect.Pointer {
		key = "*" + key
	}
	return key, nil
}
