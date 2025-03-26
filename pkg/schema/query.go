package schema

type Query map[string]any

func (q Query) init(values map[string]any) error {
	for k, v := range values {
		q[k] = v
	}
	return nil
}
