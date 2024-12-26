package html_builder_pkg

import (
	"bytes"
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/builder"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/jsonschema/v6"
)

type HtmlBuilderPkg struct {
	server  pkg.Server
	builder *builder.Builder
}

func New(server pkg.Server) *HtmlBuilderPkg {
	return NewHtmlBuilderPkg(server)
}

func NewHtmlBuilderPkg(server pkg.Server) *HtmlBuilderPkg {
	fsm := &HtmlBuilderPkg{
		server:  server,
		builder: builder.NewBuilder(),
	}
	return fsm
}

func (b *HtmlBuilderPkg) Build(ctx context.Context, schema string, tplTypeName string, opts map[string]any) string {
	if schema == "" {
		panic("Build() schema is nil")
	}
	if tplTypeName == "" {
		panic("Build() type is nil")
	}
	if opts == nil {
		opts = map[string]any{}
	}

	tType := builder.GetTplType(tplTypeName)
	if tType == builder.TplTypeNull {
		panic("Build() type is null")
	}

	schemaFile := "schema.json"
	compiler := b.server.NewSchemaCompiler()
	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(schema)))

	if err = compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}

	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}

	html, err := b.builder.CreateHTML(sch, tType, opts)
	if err != nil {
		panic(err)
	}
	return html
}
