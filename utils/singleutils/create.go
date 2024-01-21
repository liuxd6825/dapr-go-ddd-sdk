package singleutils

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/runtimeutils"
	"strings"
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

func Create[T any](fun func() T, keys ...string) T {
	//lock.Lock()
	//defer lock.Unlock()

	key := runtimeutils.GetFuncName(1)
	if len(keys) > 0 {
		key = key + "-" + strings.Join(keys, "-")
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
		}
		addItem(val)
	})
	return val.obj.(T)
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
	lock.RLock()
	defer func() {
		lock.RUnlock()
	}()
	items[item.key] = item
}
