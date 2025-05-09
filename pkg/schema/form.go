package schema

import "github.com/liuxd6825/jsonschema/v6"

type Form map[string]any

func (c Form) init(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	for k, v := range vals {
		c[k] = v
	}
	return nil
}

func NewForm() *Form {
	return &Form{}
}
