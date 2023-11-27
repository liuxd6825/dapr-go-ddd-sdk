package singleutils

import (
	"sync"
)

type onceItem struct {
	once sync.Once
	obj  any
	fun  any
	key  string
}

var lock sync.RWMutex
var items = make(map[string]*onceItem)

func Create[T any](key string, fun func() T) T {
	lock.Lock()
	defer lock.Unlock()

	val, find := items[key]
	if find {
		return val.obj.(T)
	}
	val = &onceItem{
		key: key,
		fun: fun,
	}
	val.once.Do(func() {
		if create, ok := val.fun.(func() T); ok {
			val.obj = create()
		}
		items[val.key] = val
	})
	return val.obj.(T)
}
