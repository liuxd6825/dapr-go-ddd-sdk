package ddd_neo4j

type RelationshipEntity interface {
	ElementEntity

	SetType(string)
	GetType() string

	SetStartId(int64)
	GetStartId() int64

	SetEndId(int64)
	GetEndId() int64
}

func NewRelationshipEntity() RelationshipEntity {
	return &BaseRelationship{}
}

type BaseRelationship struct {
	BaseElement
	Type    string
	StartId int64
	EndId   int64
}

func (b *BaseRelationship) SetType(s string) {
	b.Type = s
}

func (b *BaseRelationship) GetType() string {
	return b.Type
}

func (b *BaseRelationship) SetStartId(i int64) {
	b.StartId = i
}

func (b *BaseRelationship) GetStartId() int64 {
	return b.StartId
}

func (b *BaseRelationship) SetEndId(i int64) {
	b.EndId = i
}

func (b *BaseRelationship) GetEndId() int64 {
	return b.EndId
}
