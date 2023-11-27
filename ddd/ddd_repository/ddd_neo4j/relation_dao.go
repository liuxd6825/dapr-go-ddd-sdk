package ddd_neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func NewRelationDao[T Relation](driver neo4j.Driver, cypher Cypher, opts ...*Options[T]) *Dao[T] {
	dao := &Dao[T]{}
	dao.init(driver, cypher, opts...)
	return dao
}
