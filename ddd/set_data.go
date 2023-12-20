package ddd

type DataStatue int

type SetDataItem struct {
	data   Entity
	statue DataStatue
}

type SetData[T Entity] struct {
	items []*SetDataItem
}

const (
	DataStatueCreate DataStatue = iota
	DataStatueUpdate
	DataStatueDelete
	DataStatueCreateOrUpdate
)

func NewSetData[T Entity]() *SetData[T] {
	return &SetData[T]{
		items: []*SetDataItem{},
	}
}

func (d *SetData[T]) AddItems(statue DataStatue, items ...T) {
	for _, item := range items {
		i := &SetDataItem{
			statue: statue,
			data:   item,
		}
		d.items = append(d.items, i)
	}
}
func (d *SetData[T]) GetCreateList() []T {
	return d.getStatueList(DataStatueCreate)
}

func (d *SetData[T]) GetCreateOrUpdateList() []T {
	return d.getStatueList(DataStatueCreateOrUpdate)
}

func (d *SetData[T]) GetUpdateList() []T {
	return d.getStatueList(DataStatueUpdate)
}

func (d *SetData[T]) GetDeleteList() []T {
	return d.getStatueList(DataStatueDelete)
}

func (d *SetData[T]) Items() []*SetDataItem {
	return d.items
}

func (d *SetData[T]) getStatueList(statue DataStatue) []T {
	var list []T
	for _, item := range d.items {
		if item.statue == statue {
			list = append(list, item.data.(T))
		}
	}
	return list
}

func (d *SetDataItem) Data() Entity {
	return d.data
}

func (d *SetDataItem) Statue() DataStatue {
	return d.statue
}

func (t DataStatue) ToString() string {
	switch t {
	case DataStatueCreate:
		return "create"
	case DataStatueUpdate:
		return "update"
	case DataStatueDelete:
		return "delete"
	case DataStatueCreateOrUpdate:
		return "createOrUpdate"
	}
	return "none"
}
