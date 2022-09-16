package ddd_neo4j

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type Node interface {
	ddd.Entity
	GetNid() int64
	SetNid(int2 int64)

	GetId() string
	SetId(string)

	GetTenantId() string
	SetTenantId(string)

	GetLabels() []string
	SetLabels([]string)
}

type BaseNode struct {
	Id       string   `json:"id"`
	TenantId string   `json:"tenantId"`
	Nid      int64    `json:"-"`
	Labels   []string `json:"labels"`
}

func newNode() Node {
	return &BaseNode{}
}

func (b *BaseNode) GetId() string {
	return b.Id
}

func (b *BaseNode) SetId(s string) {
	b.Id = s
}

func (b *BaseNode) GetNid() int64 {
	return b.Nid
}

func (b *BaseNode) SetNid(int2 int64) {
	b.Nid = int2
}

func (b *BaseNode) GetTenantId() string {
	return b.TenantId
}

func (b *BaseNode) SetTenantId(s string) {
	b.TenantId = s
}

func (b *BaseNode) GetLabels() []string {
	return b.Labels
}

func (b *BaseNode) SetLabels(v []string) {
	b.Labels = v
}
