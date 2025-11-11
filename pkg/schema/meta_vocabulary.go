package schema

import (
	"context"
	"log"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/jsonschema/v6"
)

func NewMetaVocabulary() *jsonschema.Vocabulary {
	url := "https://extensions.com/schemas/metadata"
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(url, schema); err != nil {
		log.Fatal(err)
	}
	sch, err := c.Compile(url)
	if err != nil {
		log.Fatal(err)
	}

	return &jsonschema.Vocabulary{
		URL:    url,
		Schema: sch,
		Subschemas: []jsonschema.SchemaPath{
			{jsonschema.Prop(META_TAG_NAME), jsonschema.AllProp{}, jsonschema.AllProp{}},
		},
		Compile: metaCompile,
	}
}

func metaCompile(ctx *jsonschema.CompilerContext, obj map[string]any) (ext jsonschema.SchemaExt, err error) {
	logCtx := context.Background()
	sch := ctx.GetSchema()
	metaType := ""
	if sch == nil {
		logs.Errorfmt(logCtx, "jsonschema.metaCompile() {ctx.GetSchema() nil}")
		return nil, errors.New(" ctx.GetSchema() is nil")
	}
	location := sch.Location
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("%s %s metatype=%s", location, r, metaType)
		}
	}()

	var metaMap map[string]any
	v, _ := obj[META_TAG_NAME]
	if v != nil {
		metaMap = v.(map[string]any)
	}
	meta := NewMetaExtension()

	if err = meta.InitGraph(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitDBField(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitDBTable(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitForm(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitColumn(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitQuery(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitLang(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitParam(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitDDD(ctx, metaMap); err != nil {
		return nil, err
	}

	if err = meta.InitConvert(ctx, metaMap); err != nil {
		return nil, err
	}
	if err = meta.InitAttributes(ctx, metaMap); err != nil {
		return nil, err
	}

	return meta, err
}

func getMapItem(metaMap map[string]any, keyName string) map[string]any {
	if metaMap == nil {
		return nil
	}
	v, ok := metaMap[keyName]
	if !ok {
		return nil
	}
	keyVal, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	return keyVal
}
