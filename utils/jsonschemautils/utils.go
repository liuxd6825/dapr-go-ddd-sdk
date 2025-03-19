package jsonschemautils

import (
	"github.com/liuxd6825/jsonschema/v6"
)

func NewCompiler() *jsonschema.Compiler {
	compiler := jsonschema.NewCompiler()
	initCompiler(compiler)
	return compiler
}

func NewJsonSchemaWidthJson(fileName string, jsonText string) *jsonschema.Schema {
	s, err := jsonschema.NewSchemaWithJson(fileName, jsonText, func(compiler *jsonschema.Compiler) {
		initCompiler(compiler)
	})
	if err != nil {
		panic(err)
	}
	return s
}

func initCompiler(compiler *jsonschema.Compiler) {
	compiler.AssertFormat()
	compiler.AssertContent()
	compiler.RegisterFormat(DateTimeFormat)
	compiler.RegisterFormat(DateFormat)
}
