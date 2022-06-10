package ddd

type AggregateId interface {
	RootId() string
	ItemIds() *[]string
	ItemCount() int
	ItemEmpty() bool
}

type aggregateId struct {
	rootId  string
	itemIds *[]string
}

func NewAggregateId(rootId string, itemIds ...string) AggregateId {
	return &aggregateId{
		rootId:  rootId,
		itemIds: &itemIds,
	}
}

func NewAggregateIds(rootId string, itemIds *[]string) AggregateId {
	return &aggregateId{
		rootId:  rootId,
		itemIds: itemIds,
	}
}

func (a *aggregateId) RootId() string {
	return a.rootId
}

func (a *aggregateId) ItemIds() *[]string {
	return a.itemIds
}

func (a *aggregateId) ItemCount() int {
	if a.itemIds == nil {
		return 0
	}
	return len(*a.itemIds)
}

func (a *aggregateId) ItemEmpty() bool {
	if a.ItemCount() == 0 {
		return true
	}
	return false
}
