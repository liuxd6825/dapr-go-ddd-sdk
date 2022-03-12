package ddd_repository

type EntityBuilder interface {
	New() interface{}
	NewList() interface{}
}
