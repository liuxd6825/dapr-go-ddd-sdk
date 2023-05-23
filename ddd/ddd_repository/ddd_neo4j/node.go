package ddd_neo4j

type Node interface {
	GetTenantId() string
	SetTenantId(v string)

	GetId() string
	SetId(v string)

	GetNid() int64
	SetNid(int2 int64)

	GetLabels() []string
	SetLabels([]string)

	SetGraphId(v string)
	GetGraphId() string
}

type BaseNode struct {
	Id       string   `json:"id" bson:"id"`
	TenantId string   `json:"tenantId" bson:"tenant_id" gorm:"index:idx_tenant_id"`
	GraphId  string   `json:"graphId" bson:"graph_id" gorm:"index:idx_graph_id"`
	Nid      int64    `json:"-" bson:"nid"`
	Labels   []string `json:"labels" bson:"labels" gorm:"-"`
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

func (b *BaseNode) GetTenantId() string {
	return b.TenantId
}

func (b *BaseNode) SetTenantId(s string) {
	b.TenantId = s
}

func (b *BaseNode) SetGraphId(v string) {
	b.GraphId = v
}

func (b *BaseNode) GetGraphId() string {
	return b.GraphId
}

func (b *BaseNode) GetNid() int64 {
	return b.Nid
}

func (b *BaseNode) SetNid(int2 int64) {
	b.Nid = int2
}

func (b *BaseNode) GetLabels() []string {
	return b.Labels
}

func (b *BaseNode) SetLabels(v []string ) {
	b.Labels = v
}
