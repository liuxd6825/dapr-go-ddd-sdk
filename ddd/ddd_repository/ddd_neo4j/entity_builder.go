package ddd_neo4j

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
)

type NodeEntityBuilder[T any] interface {
	ddd.EntityBuilder[T]
	GetLabels(e T) []string
}

type RelationEntityBuilder[T any] interface {
	ddd.EntityBuilder[T]
	GetStartId(entity T) string
	GetEndId(entity T) string
	GetRelType(entity T) string
}

type nodeEntityBuilder[T any] struct {
	ddd.EntityBuilder[T]
}
type relationEntityBuilder[T any] struct {
	ddd.EntityBuilder[T]
}

func NewRelationEntityBuilder[T any]() RelationEntityBuilder[T] {
	base := ddd.NewAnyEntityBuilder[T]()
	return &relationEntityBuilder[T]{
		EntityBuilder: base,
	}
}

func NewNodeEntityBuilder[T any]() NodeEntityBuilder[T] {
	base := ddd.NewAnyEntityBuilder[T]()
	return &nodeEntityBuilder[T]{
		EntityBuilder: base,
	}
}

func (r *nodeEntityBuilder[T]) GetLabels(entity T) []string {
	return reflectutils.GetFieldStrings(entity, "labels")
}

func (r *relationEntityBuilder[T]) GetStartId(entity T) string {
	return reflectutils.GetFieldString(entity, "startId")
}

func (r *relationEntityBuilder[T]) SetStartId(entity T, val string) {
	reflectutils.SetFieldString(entity, "startId", val)
}

func (r *relationEntityBuilder[T]) GetEndId(entity T) string {
	return reflectutils.GetFieldString(entity, "endId")
}

func (r *relationEntityBuilder[T]) SetEndId(entity T, val string) {
	reflectutils.SetFieldString(entity, "endId", val)
}

func (r *relationEntityBuilder[T]) GetRelType(entity T) string {
	return reflectutils.GetFieldString(entity, "relType")
}

func (r *relationEntityBuilder[T]) SetRelType(entity T, val string) {
	reflectutils.SetFieldString(entity, "relType", val)
}
