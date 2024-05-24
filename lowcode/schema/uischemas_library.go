package schema

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/maputils"
	"strings"
)

type UiSchemasLibrary struct {
	items map[string]*UiSchema
}

func NewUiSchemasLibrary() *UiSchemasLibrary {
	return &UiSchemasLibrary{items: make(map[string]*UiSchema)}
}

// LoadJavaScript
//
//	@Description: 加载js代码字符串
//	@receiver s
//	@param jsCode
//	@return error
func (s *UiSchemasLibrary) LoadJavaScript(bytes []byte) error {
	items, err := s.loadCode(bytes)
	if err != nil {
		return err
	}
	s.items = items
	return nil
}

func (s *UiSchemasLibrary) Get(name string) (*UiSchema, error) {
	schema, ok := s.items[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("uischema %s not found", name))
	}
	return schema, nil
}

func (s *UiSchemasLibrary) Set(name string, uiSchema *UiSchema) {
	s.items[name] = uiSchema
}

func (s *UiSchemasLibrary) Map() map[string]*UiSchema {
	return s.items
}

func (s *UiSchemasLibrary) loadCode(bytes []byte) (map[string]*UiSchema, error) {
	str := string(bytes)
	code := strings.ReplaceAll(str, "export var", " var")
	vm := goja.New()
	_, err := vm.RunString(code)
	if err != nil {
		return nil, err
	}
	value := vm.Get("uiSchemas")
	if value == nil {
		return nil, errors.New("uiSchemas is not an object")
	}

	data, ok := value.Export().(map[string]any)
	if !ok {
		return nil, errors.New("uiSchemas is not an map[string]any")
	}
	items := map[string]*UiSchema{}
	for key, item := range data {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, errors.New(fmt.Sprintf("uiSchema %s is not an object", key))
		}
		var uiSchema UiSchema
		if err := maputils.Decode(m, &uiSchema); err != nil {
			return nil, err
		}
		items[key] = &uiSchema
	}
	return items, nil
}
