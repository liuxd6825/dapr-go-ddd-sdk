package html_pkg

import (
	"bytes"
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/fs_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/html/builder"
	"github.com/liuxd6825/jsonschema/v6"
)

type HtmlPkg struct {
	server  pkg.Server
	builder *builder.Builder
	*fs_pkg.FsPkg
}

func New(server pkg.Server) *HtmlPkg {
	return NewHtmlPkg(server)
}

func NewHtmlPkg(server pkg.Server) *HtmlPkg {
	fsPkg, err := fs_pkg.NewFsPkg(server.GetEnvConfig(), "web")
	if err != nil {
		panic(err)
	}
	fsm := &HtmlPkg{
		server:  server,
		builder: builder.NewBuilder(),
		FsPkg:   fsPkg,
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
