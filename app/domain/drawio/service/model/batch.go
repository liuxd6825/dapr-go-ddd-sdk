package model

type SaveBatch struct {
	Nodes     *Batch[*Node]
	Relations *Batch[*Relation]
}

type Batch[T Item] struct {
	Updates map[string]T
	Removes map[string]T
	Creates map[string]T
}

type Item interface {
	GetId() string
}

func NewSaveBatch() *SaveBatch {
	return &SaveBatch{
		Nodes:     NewBatch[*Node](),
		Relations: NewBatch[*Relation](),
	}
}

func NewBatch[T Item]() *Batch[T] {
	return &Batch[T]{
		Updates: make(map[string]T, 0),
		Removes: make(map[string]T, 0),
		Creates: make(map[string]T, 0),
	}
}

func (b *Batch[T]) AddUpdate(updates ...T) {
	for _, update := range updates {
		b.Updates[update.GetId()] = update
	}
}

func (b *Batch[T]) AddRemove(remove ...T) {
	for _, item := range remove {
		b.Removes[item.GetId()] = item
	}
}

func (b *Batch[T]) AddCreate(insert ...T) {
	for _, item := range insert {
		b.Creates[item.GetId()] = item
	}
}
