package form

import (
	"bytes"
	"github.com/liuxd6825/jsonschema/v6"
	"os"
	"testing"
)

// GenerateHTML generates HTML code for ui5-form based on JsonSchema
func TestGenerateHTML(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path + "/testfile/human.json")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := "schema.json"
	compiler := jsonschema.NewCompiler()

	if err := compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}

	builder := NewBuilder()
	html, err := builder.CreateHTML(sch)
	if err != nil {
		panic(err)
	}
	t.Log(html)
}
