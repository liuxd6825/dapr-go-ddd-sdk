package ddd_repository

import "github.com/dapr/dapr-go-ddd-sdk/ddd"

/*
type EntityBuilder[T any] interface {
	NewOne() T
	NewList() T
}
*/

type EntityBuilder[T ddd.Entity] struct {
	newOneFunc  func() T
	newListFunc func() *[]T
}

func (e *EntityBuilder[T]) NewOne() T {
	return e.newOneFunc()
}

func (e *EntityBuilder[T]) NewList() *[]T {
	return e.newListFunc()
}

func NewEntityBuilder[T ddd.Entity](newOne func() T, newList func() *[]T) *EntityBuilder[T] {
	return &EntityBuilder[T]{
		newOneFunc:  newOne,
		newListFunc: newList,
	}
}
