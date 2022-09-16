package ddd_neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type RelationDao[T Relation] struct {
	BaseDao[T]
}

func NewRelationDao[T Relation](driver neo4j.Driver, cypher Cypher, opts ...*Options[T]) *RelationDao[T] {
	dao := &RelationDao[T]{}
	dao.init(driver, cypher, opts...)
	return dao
}
