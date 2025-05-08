package ddd_neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewRelationDao[T Relation](driver neo4j.DriverWithContext, cypher Cypher, opts ...*Options[T]) *Dao[T] {
	dao := &Dao[T]{}
	dao.init(driver, cypher, opts...)
	return dao
}
