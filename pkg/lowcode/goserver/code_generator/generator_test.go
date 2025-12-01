package code_generator

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/liuxd6825/jsonschema/v6"
)

//go:embed human.json
var humanJsonSchema []byte

func Test_Generate(t *testing.T) {
	generator := NewGenerator()
	sch, err := jsonschema.NewSchemaWithBytes("human.json", humanJsonSchema)
	if err != nil {
		t.Fatal(err)
		return
	}
	// 生成代码
	code, err := generator.GeneratorModel(sch, &Options{
		PackageName: "model",
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(string(code))

}
