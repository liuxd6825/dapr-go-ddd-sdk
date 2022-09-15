package ddd_neo4j

import "github.com/liuxd6825/dapr-go-ddd-sdk/ddd"

type Relation interface {
	ddd.Entity

	SetNid(int64)
	GetNid() int64

	SetSid(int64)
	GetSid() int64

	SetEid(v int64)
	GetEid() int64

	SetType(string)
	GetType() string

	SetStartId(string)
	GetStartId() string

	SetEndId(v string)
	GetEndId() string
}

type BaseRelation struct {
	Id       string `json:"id"`
	TenantId string `json:"tenantId"`
	Nid      int64  `json:"-"`
	Sid      int64  `json:"-"`
	Eid      int64  `json:"-"`
	Type     string `json:"type"`
	StartId  string `json:"startId"`
	EndId    string `json:"endId"`
}

func newRelation() Relation {
	return &BaseRelation{}
}

func (b *BaseRelation) SetNid(s int64) {
	b.Nid = s
}

func (b *BaseRelation) GetNid() int64 {
	return b.Nid
}

func (b *BaseRelation) GetTenantId() string {
	return b.TenantId
}

func (b *BaseRelation) SetTenantId(v string) {
	b.TenantId = v
}

func (b *BaseRelation) GetId() string {
	return b.Id
}

func (b *BaseRelation) SetId(v string) {
	b.Id = v
}

func (b *BaseRelation) SetType(s string) {
	b.Type = s
}

func (b *BaseRelation) GetType() string {
	return b.Type
}

func (b *BaseRelation) SetSid(s int64) {
	b.Sid = s
}

func (b *BaseRelation) GetSid() int64 {
	return b.Sid
}

func (b *BaseRelation) SetEid(s int64) {
	b.Eid = s
}

func (b *BaseRelation) GetEid() int64 {
	return b.Eid
}

func (b *BaseRelation) SetStartId(s string) {
	b.StartId = s
}

func (b *BaseRelation) GetStartId() string {
	return b.StartId
}

func (b *BaseRelation) SetEndId(v string) {
	b.EndId = v
}

func (b *BaseRelation) GetEndId() string {
	return b.EndId
}
