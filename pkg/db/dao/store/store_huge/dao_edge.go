package store_huge

import "github.com/apache/hugegraph-toolchain/hugegraph-client-go"

type EdgeDao[T IEdge] struct {
	*BaseDao[T]
}

func NewEdgeDao[T IEdge](client *hugegraph.CommonClient, config *Config[T]) *EdgeDao[T] {
	return &EdgeDao[T]{
		BaseDao: NewBaseDao[T](client, config),
	}
}
