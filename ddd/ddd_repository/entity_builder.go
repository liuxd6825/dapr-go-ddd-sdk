package ddd_repository

type EntityBuilder interface {
	New() interface{}
	NewList() interface{}
}

type NewFunc func() interface{}

func NewEntityBuilder(newOne NewFunc, newList NewFunc) EntityBuilder {
	return &entityBuilder{
		newOne:  newOne,
		newList: newList,
	}
}

type entityBuilder struct {
	newOne  NewFunc
	newList NewFunc
}

func (e *entityBuilder) New() interface{} {
	return e.newOne()
}

func (e *entityBuilder) NewList() interface{} {
	return e.newList()
}
