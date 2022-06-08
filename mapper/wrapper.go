package mapper

import (
	"reflect"
)

type TypeWrapper interface {
	IsType(value reflect.Value) bool
	SetValue(value reflect.Value, toFieldInfo reflect.Value) (bool, error)
	SetNext(m TypeWrapper)
}

type BaseTypeWrapper struct {
	next TypeWrapper
}

func (bm *BaseTypeWrapper) SetNext(m TypeWrapper) {
	bm.next = m
}
