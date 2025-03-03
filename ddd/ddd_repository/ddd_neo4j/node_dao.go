package ddd_neo4j

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewNodeDao[T Node](driver neo4j.DriverWithContext, cypher Cypher, eb ddd.EntityBuilder[T], opts ...*Options[T]) *Dao[T] {
	dao := &Dao[T]{}
	dao.init(driver, cypher, eb, opts...)
	return dao
}
