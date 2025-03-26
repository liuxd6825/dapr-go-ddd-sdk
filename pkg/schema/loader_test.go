package schema

import (
	"fmt"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/spf13/afero"
	"github.com/xeipuuv/gojsonschema"
	"testing"
)

func Test_JSONLoader(t *testing.T) {
	// 定义主 Schema 的路径
	mainSchemaPath := "./testfile/person.json"
	fs := afero.NewOsFs()
	// 创建 SchemaLoader 实例
	compiler := jsonschema.NewCompiler()

	// 验证数据
	data := `
	{
		"name": "John Doe",
		"address": {
			"street": "123 Main St",
			"city": "Springfield",
			"zipcode": "12345"
		}
	}`
	documentLoader := gojsonschema.NewStringLoader(data)
	schema, err := compiler.Compile(mainSchemaPath)
	if err != nil {
		t.Fatal(err)
	}

	// 执行验证
	err = schema.Validate(documentLoader)
	if err != nil {
		fmt.Println("Validation error:", err)
		return
	}
}
