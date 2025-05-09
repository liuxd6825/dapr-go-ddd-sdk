package schema

import "github.com/liuxd6825/jsonschema/v6"

type Query map[string]any

func NewQuery() *Query {
	return &Query{}
}

func (q Query) init(ctx *jsonschema.CompilerContext, values map[string]any) error {
	for k, v := range values {
		q[k] = v
	}
	return nil
}
