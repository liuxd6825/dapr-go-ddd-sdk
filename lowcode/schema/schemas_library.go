package schema

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"strings"
)

type SchemasLibrary struct {
	schemas map[string]*Schema
}

func NewSchemasLibrary() *SchemasLibrary {
	return &SchemasLibrary{schemas: make(map[string]*Schema)}
}

func (s *SchemasLibrary) LoadJavaScript(bytes []byte) error {
	items, err := s.loadJsCode(bytes)
	if err != nil {
		return err
	}
	s.schemas = items
	return nil
}

func (s *SchemasLibrary) loadJsCode(bytes []byte) (map[string]*Schema, error) {
	code := strings.ReplaceAll(string(bytes), "export var", " var")
	vm := goja.New()
	_, err := vm.RunString(code)
	if err != nil {
		return nil, err
	}
	value := vm.Get("schemas")
	if value == nil {
		return nil, errors.New("schemas.SchemasLibrary is not an object")
	}

	data, ok := value.Export().(map[string]any)
	if !ok {
		return nil, errors.New("schemas.SchemasLibrary is not an map[string]any")
	}
	items := map[string]*Schema{}
	for key, item := range data {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, errors.New(fmt.Sprintf("schema %s is not an object", key))
		}
		var schema Schema
		if err := maputils.Decode(m, &schema); err != nil {
			return nil, err
		}
		items[key] = &schema
	}
	return items, nil
}

func (s *SchemasLibrary) Get(name string) (*Schema, error) {
	schema, ok := s.schemas[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("schema %s not found", name))
	}
	return schema, nil
}

func (s *SchemasLibrary) Set(name string, schema *Schema) {
	s.schemas[name] = schema
}

func (s *SchemasLibrary) Map() map[string]*Schema {
	return s.schemas
}
