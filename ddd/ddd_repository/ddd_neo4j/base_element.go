package ddd_neo4j

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type ElementEntity interface {
	ddd.Entity
	SetNid(int2 int64)
	GetNid() int64

	SetLabels([]string)
	GetLabels() []string

	GetId() string
	SetId(string)

	SetTenantId(string)
	GetTenantId() string
}

type BaseElement struct {
	Id       string   `json:"id"`
	TenantId string   `json:"tenantId"`
	Nid      int64    `json:"-"`
	Labels   []string `json:"-"`
}

func (b *BaseElement) SetNid(int2 int64) {
	b.Nid = int2
}

func (b *BaseElement) GetNid() int64 {
	return b.Nid
}

func (b *BaseElement) SetLabels(strings []string) {
	b.Labels = strings
}

func (b *BaseElement) GetLabels() []string {
	return b.Labels
}

func (b *BaseElement) SetTenantId(s string) {
	b.TenantId = s
}

func (b *BaseElement) GetTenantId() string {
	return b.TenantId
}

func (b *BaseElement) SetId(s string) {
	b.Id = s
}

func (b *BaseElement) GetId() string {
	return b.Id
}

func newElementEntity() ElementEntity {
	return &BaseElement{}
}
