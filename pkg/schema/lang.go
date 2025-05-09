package schema

import (
	"fmt"
	"github.com/liuxd6825/jsonschema/v6"
)

type Lang map[string]string

func NewLang() *Lang {
	return &Lang{}
}

func (l Lang) init(ctx *jsonschema.CompilerContext, values map[string]any) error {
	for k, v := range values {
		l[k] = fmt.Sprintf("%s", v)
	}
	return nil
}
