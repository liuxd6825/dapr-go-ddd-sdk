package model

type SaveBatch struct {
	Nodes     *ItemBatch[*Node]
	Relations *ItemBatch[*Relation]
}

type ItemBatch[T Item] struct {
	Updates map[string]T
	Removes map[string]T
	Creates map[string]T
}

type Item interface {
	GetId() string
}

func NewSaveBatch() *SaveBatch {
	return &SaveBatch{
		Nodes:     NewItemBatch[*Node](),
		Relations: NewItemBatch[*Relation](),
	}
}

func NewItemBatch[T Item]() *ItemBatch[T] {
	return &ItemBatch[T]{
		Updates: make(map[string]T, 0),
		Removes: make(map[string]T, 0),
		Creates: make(map[string]T, 0),
	}
}

func (b *ItemBatch[T]) AddUpdate(updates ...T) {
	for _, update := range updates {
		b.Updates[update.GetId()] = update
	}
}

func (b *ItemBatch[T]) AddRemove(remove ...T) {
	for _, item := range remove {
		b.Removes[item.GetId()] = item
	}
}

func (b *ItemBatch[T]) AddCreate(insert ...T) {
	for _, item := range insert {
		b.Creates[item.GetId()] = item
	}
}
