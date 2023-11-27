package ddd_neo4j

import (
	"context"
	"github.com/google/uuid"
	"testing"
)

type node struct {
	tenantId string
	id       string
	nid      int64
	labels   []string
	graphId  string
	Name     string
	Title    string
}

func TestNodeCypher_InsertOrUpdate(t *testing.T) {
	ctx := context.Background()
	node := &node{
		tenantId: "test",
		id:       uuid.NewString(),
		graphId:  "test",
		Name:     "testName",
		Title:    "testTitle",
	}
	nodeCypher := NewNodeCypher("test")
	cr, err := nodeCypher.InsertOrUpdate(ctx, node)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cr.Cypher())
	}
}

func (n *node) GetTenantId() string {
	return n.tenantId
}

func (n *node) SetTenantId(v string) {
	n.tenantId = v
}

func (n *node) GetId() string {
	return n.id
}

func (n *node) SetId(v string) {
	n.id = v
}

func (n *node) GetNid() int64 {
	return n.nid
}

func (n *node) SetNid(int2 int64) {
	n.SetNid(int2)
}

func (n *node) GetLabels() []string {
	return n.labels
}

func (n *node) SetLabels(strings []string) {
	n.labels = strings
}

func (n *node) SetGraphId(v string) {
	n.graphId = v
}

func (n *node) GetGraphId() string {
	return n.graphId
}
