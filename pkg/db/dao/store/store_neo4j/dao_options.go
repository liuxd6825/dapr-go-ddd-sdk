package store_neo4j

type Options[T interface{}] struct {
	newOne  func() T
	newList func() []T
}

func NewOptions[T interface{}](opts ...*Options[T]) *Options[T] {
	n := &Options[T]{}
	for _, o := range opts {
		if o.newList != nil {
			n.newList = o.newList
		}
		if o.newOne != nil {
			n.newOne = o.newOne
		}
	}
	return n
}
