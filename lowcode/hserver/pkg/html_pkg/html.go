package html_pkg

import (
	"bytes"
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/html/builder"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/fspkg"
	"github.com/liuxd6825/jsonschema/v6"
)

type HtmlPkg struct {
	server  element.Server
	builder *builder.Builder
	fspkg.IFsPkg
}

func New(server element.Server) *HtmlPkg {
	return NewHtmlPkg(server)
}

func NewHtmlPkg(server element.Server) *HtmlPkg {
	fsm := &HtmlPkg{
		server:  server,
		builder: builder.NewBuilder(),
	}
	return fsm
}

func (b *HtmlPkg) Build(ctx context.Context, schema string, typeName string, opts map[string]any) string {
	if schema == "" {
		panic("Build() schema is nil")
	}
	if typeName == "" {
		panic("Build() typeName is nil")
	}
	if opts == nil {
		opts = map[string]any{}
	}

	tType := builder.NewTemplateType(typeName)
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
