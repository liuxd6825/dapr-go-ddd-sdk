package ddd_neo4j

type Relation interface {
	GetTenantId() string
	SetTenantId(v string)

	GetId() string
	SetId(v string)

	SetNid(int64)
	GetNid() int64

	SetSid(int64)
	GetSid() int64

	SetEid(v int64)
	GetEid() int64

	SetRelType(string)
	GetRelType() string

	SetStartId(string)
	GetStartId() string

	SetEndId(v string)
	GetEndId() string

	SetGraphId(v string)
	GetGraphId() string

	SetProperties(v map[string]interface{})
	GetProperties() map[string]interface{}
}

type BaseRelation struct {
	Id         string                 `json:"id" bson:"id"`
	TenantId   string                 `json:"tenantId" bson:"tenant_id" gorm:"index:idx_tenant_id"`
	GraphId    string                 `json:"graphId" bson:"graph_id" gorm:"index:idx_graph_id"`
	Nid        int64                  `json:"-" bson:"nid" gorm:"-"`
	Sid        int64                  `json:"-" bson:"sid" gorm:"-"`
	Eid        int64                  `json:"-" bson:"eid" gorm:"-"`
	RelType    string                 `json:"relType" bson:"rel_type" gorm:"index:idx_rel_type"`
	StartId    string                 `json:"startId" bson:"start_id" gorm:"index:idx_start_id"`
	EndId      string                 `json:"endId" bson:"end_id" gorm:"index:idx_end_id"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

func newRelation() Relation {
	return &BaseRelation{}
}

func (b *BaseRelation) GetId() string {
	return b.Id
}

func (b *BaseRelation) SetId(v string) {
	b.Id = v
}

func (b *BaseRelation) GetTenantId() string {
	return b.TenantId
}

func (b *BaseRelation) SetTenantId(v string) {
	b.TenantId = v
}

func (b *BaseRelation) SetGraphId(v string) {
	b.GraphId = v
}

func (b *BaseRelation) GetGraphId() string {
	return b.GraphId
}

func (b *BaseRelation) SetNid(s int64) {
	b.Nid = s
}

func (b *BaseRelation) GetNid() int64 {
	return b.Nid
}

func (b *BaseRelation) SetRelType(s string) {
	b.RelType = s
}

func (b *BaseRelation) GetRelType() string {
	return b.RelType
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

func (b *BaseRelation) SetProperties(v map[string]interface{}) {
	b.Properties = v
}

func (b *BaseRelation) GetProperties() map[string]interface{} {
	return b.Properties
}
