package ddd_repository

type EntityBuilder interface {
	NewOne() interface{}
	NewList() interface{}
}

type entityBuilder struct {
	newOneFunc  func() interface{}
	newListFunc func() interface{}
}

func (e *entityBuilder) NewOne() interface{} {
	return e.newOneFunc()
}

func (e *entityBuilder) NewList() interface{} {
	return e.newListFunc()
}

func NewEntityBuilder(newOne func() interface{}, newList func() interface{}) EntityBuilder {
	return &entityBuilder{
		newOneFunc:  newOne,
		newListFunc: newList,
	}
}
