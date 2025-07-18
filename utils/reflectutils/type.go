package reflectutils

import (
	"reflect"
	"sync"
	"unsafe"
)

// TypeDetailsCache 是一个线程安全的缓存，用于存储类型信息
type TypeDetailsCache struct {
	cache map[TypeKey]typeDetails
	mu    sync.RWMutex
}

type typeDetails struct {
	pkgPath string
	name    string
	fields  map[string]reflect.StructField
}

var typeCache = &TypeDetailsCache{
	cache: make(map[TypeKey]typeDetails),
}

// GetTypeKey 获取类型唯一标识（不依赖reflect.TypeOf）
func GetTypeKey(obj interface{}) uintptr {
	// 使用unsafe获取对象类型标识
	typePtr := (*[2]uintptr)(unsafe.Pointer(&obj))[1]
	return typePtr
}

// GetTypeDetails 返回对象的类型详细信息，优先从缓存获取
func GetTypeDetails(obj interface{}) (pkgPath, typeName string, fields map[string]reflect.StructField) {
	// 第一步：尝试从缓存获取类型信息（无需先获取类型）
	typeKey, cached := getCachedTypeDetails(obj)
	if cached != nil {
		return cached.pkgPath, cached.name, cached.fields
	}

	// 第二步：缓存未命中时再进行反射操作
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 如果未找到，使用反射获取信息并缓存
	details := typeDetails{
		pkgPath: t.PkgPath(),
		name:    t.Name(),
		fields:  make(map[string]reflect.StructField),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		details.fields[field.Name] = field
	}
	typeCache.mu.Lock()
	typeCache.cache[typeKey] = details
	typeCache.mu.Unlock()

	return details.pkgPath, details.name, details.fields
}

// getCachedTypeDetails 从缓存中获取类型信息
func getCachedTypeDetails(obj interface{}) (TypeKey, *typeDetails) {
	// 使用unsafe获取类型标识作为缓存key
	typeKey := getTypeKey(obj)

	typeCache.mu.RLock()
	defer typeCache.mu.RUnlock()

	if details, ok := typeCache.cache[typeKey]; ok {
		return typeKey, &details
	}
	return typeKey, nil
}

// typeKey 作为缓存键的类型
type TypeKey uintptr

// getTypeKey 获取类型的唯一标识
func getTypeKey(obj interface{}) TypeKey {
	// 获取对象类型指针
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 使用reflect.Type的指针作为key
	return TypeKey(reflect.ValueOf(t).Pointer())
}
