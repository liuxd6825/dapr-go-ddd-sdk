package model

type SaveBatch struct {
	Nodes     *Batch[*Node]
	Relations *Batch[*Relation]
}

type Batch[T any] struct {
	Updates []T
	Removes []T
	Inserts []T
}

func NewSaveBatch() *SaveBatch {
	return &SaveBatch{
		Nodes:     NewBatch[*Node](),
		Relations: NewBatch[*Relation](),
	}
}

func NewBatch[T any]() *Batch[T] {
	return &Batch[T]{
		Updates: make([]T, 0),
		Removes: make([]T, 0),
		Inserts: make([]T, 0),
	}
}

func (b *Batch[T]) Update(updates ...T) {
	b.Updates = append(b.Updates, updates...)
}

func (b *Batch[T]) Remove(remove ...T) {
	b.Removes = append(b.Removes, remove...)
}

func (b *Batch[T]) Insert(insert ...T) {
	b.Inserts = append(b.Inserts, insert...)
}
