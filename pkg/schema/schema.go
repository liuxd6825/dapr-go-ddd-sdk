package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/jsonschema/v6"
)

type Schema jsonschema.Schema

func NewCompiler() *jsonschema.Compiler {
	compiler := jsonschema.NewCompiler()
	initCompiler(compiler)
	return compiler
}

func NewJsonSchemaWithJson(fileName string, jsonText string) *jsonschema.Schema {
	sch, err := jsonschema.NewSchemaWithJson(fileName, jsonText, func(compiler *jsonschema.Compiler) {
		initCompiler(compiler)
	})
	if err != nil {
		panic(err)
	}
	return sch
}

func NewJsonSchemaWithBytes(fileName string, fileBytes []byte) *jsonschema.Schema {
	sch, err := jsonschema.NewSchemaWithBytes(fileName, fileBytes, func(compiler *jsonschema.Compiler) {
		initCompiler(compiler)
	})
	if err != nil {
		panic(err)
	}
	return sch
}

func NewJsonSchemaWithStruct(fileName string, obj any) *jsonschema.Schema {
	sch, err := jsonschema.NewSchemaWithStruct(fileName, obj, func(compiler *jsonschema.Compiler) {
		initCompiler(compiler)
	})
	if err != nil {
		panic(err)
	}
	return sch
}

func initCompiler(compiler *jsonschema.Compiler) {
	compiler.AssertFormat()
	compiler.AssertContent()
	compiler.AssertVocabs()

	compiler.RegisterFormat(DateTimeFormat)
	compiler.RegisterFormat(DateFormat)
	compiler.RegisterVocabulary(NewMetaVocabulary())

	if env.GetEnv() != nil && env.GetEnv().Fsm != nil {
		compiler.UseLoader(NewJSONLoader())
	}

}
