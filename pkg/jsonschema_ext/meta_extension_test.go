package jsonschema_ext

import (
	"bytes"
	"github.com/liuxd6825/jsonschema/v6"
	"os"
	"testing"
)

func TestSchema_MetaExtension(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path + "/testdata/example.json")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := "schema.json"
	compiler := jsonschema.NewCompiler()
	compiler.AssertVocabs()
	compiler.RegisterVocabulary(NewMetaVocabulary())

	if err := compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}

	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}

	props := sch.GetSortProperties()
	for _, prop := range props {
		println(prop.Name())
	}
	t.Log(props)
}
