package store_neo4j

import (
	"fmt"

	store2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/reflectutils"
)

type NodeEntityBuilder[T any] interface {
	store2.EntityBuilder[T]
	GetLabels(e T) []string
}

type RelationEntityBuilder[T any] interface {
	store2.EntityBuilder[T]
	GetStartId(entity T) string
	GetEndId(entity T) string
	GetRelType(entity T) string
}

type nodeEntityBuilder[T any] struct {
	*store2.AnyEntityBuilder[T]
}
type relationEntityBuilder[T any] struct {
	*store2.AnyEntityBuilder[T]
}

func NewRelationEntityBuilder[T any](sch *store2.DBSchema) RelationEntityBuilder[T] {
	base := store2.NewAnyEntityBuilder[T](sch)
	return &relationEntityBuilder[T]{
		AnyEntityBuilder: base,
	}
}

func NewNodeEntityBuilder[T any](sch *store2.DBSchema) NodeEntityBuilder[T] {
	base := store2.NewAnyEntityBuilder[T](sch)
	return &nodeEntityBuilder[T]{
		AnyEntityBuilder: base,
	}
}

func (r *nodeEntityBuilder[T]) GetLabels(entity T) []string {
	var labels []string
	fields := r.GetDBSchema().GetNodeLabelFields()
	for _, field := range fields {
		if field != nil && field.NodeLabel {
			label := reflectutils.GetFieldString(entity, field.Name)
			if field.NodeLabelFormat != "" {
				labels = append(labels, fmt.Sprintf(field.NodeLabelFormat, label))
			} else {
				labels = append(labels, label)
			}
		}
	}
	return labels
}

func (r *relationEntityBuilder[T]) GetStartId(entity T) string {
	field := r.GetDBSchema().GetRelStartIdField()
	if field != nil {
		propName := r.GetPropertyName(field)
		return reflectutils.GetFieldString(entity, propName)
	}
	return ""
}

func (r *relationEntityBuilder[T]) SetStartId(entity T, val string) {
	field := r.GetDBSchema().GetRelStartIdField()
	if field != nil {
		propName := r.GetPropertyName(field)
		reflectutils.SetFieldString(entity, propName, val)
	}
}

func (r *relationEntityBuilder[T]) GetEndId(entity T) string {
	field := r.GetDBSchema().GetRelEndIdField()
	if field != nil {
		propName := r.GetPropertyName(field)
		return reflectutils.GetFieldString(entity, propName)
	}
	return ""
}

func (r *relationEntityBuilder[T]) SetEndId(entity T, val string) {
	field := r.GetDBSchema().GetRelEndIdField()
	if field != nil {
		propName := r.GetPropertyName(field)
		reflectutils.SetFieldString(entity, propName, val)
	}
}

func (r *relationEntityBuilder[T]) GetRelType(entity T) string {
	field := r.GetDBSchema().GetRelTypeField()
	if field != nil {
		propName := r.GetPropertyName(field)
		return reflectutils.GetFieldString(entity, propName)
	}
	return ""
}

func (r *relationEntityBuilder[T]) SetRelType(entity T, val string) {
	field := r.GetDBSchema().GetRelTypeField()
	if field != nil {
		propName := r.GetPropertyName(field)
		reflectutils.SetFieldString(entity, propName, val)
	}
	return
}
