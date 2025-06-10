package store_neo4j

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

type TNode struct {
	TenantId string
	Id       string
	Nid      int64
	Labels   []string
	GraphId  string
	Name     string
	Title    string
}

func TestNodeCypher_InsertOrUpdate(t *testing.T) {
	ctx := context.Background()
	node := &TNode{
		TenantId: "test",
		Id:       uuid.NewString(),
		GraphId:  "test",
		Name:     "testName",
		Title:    "testTitle",
	}

	eb := NewNodeEntityBuilder[*TNode](nil)
	nodeCypher := NewNodeCypher(&Config[*TNode]{
		EntityBuilder: eb,
		DBSchema:      nil,
		Labels:        make([]string, 0),
	})
	cr, err := nodeCypher.InsertOrUpdate(ctx, node)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cr.Cypher())
	}
}
