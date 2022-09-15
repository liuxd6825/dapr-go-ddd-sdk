package ddd_neo4j

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type NodeDao[T Node] struct {
	BaseDao[T]
}

func NewNodeDao[T Node](driver neo4j.Driver, cypher Cypher) *NodeDao[T] {
	dao := &NodeDao[T]{}
	dao.init(driver, cypher)
	return dao
}
