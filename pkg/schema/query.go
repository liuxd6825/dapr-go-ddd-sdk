package schema

type Query map[string]any

func NewQuery() *Query {
	return &Query{}
}

func (q Query) init(values map[string]any) error {
	for k, v := range values {
		q[k] = v
	}
	return nil
}
