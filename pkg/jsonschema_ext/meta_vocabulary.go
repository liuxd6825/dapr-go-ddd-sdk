package jsonschema_ext

import (
	"github.com/liuxd6825/jsonschema/v6"
	"log"
	"strings"
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

func metaCompile(ctx *jsonschema.CompilerContext, obj map[string]any) (jsonschema.SchemaExt, error) {
	v, ok := obj[META_TAG_NAME]
	if !ok {
		return nil, nil
	}
	d, ok := v.(map[string]any)
	if !ok {
		return nil, nil
	}
	var err error
	meta := NewMetaExtension()
	for key, value := range d {
		switch key {
		case "dbField":
			vals, ok := value.(map[string]any)
			if ok {
				err = meta.InitDBField(vals)
			}
		case "dbTable":
			vals, ok := value.(map[string]any)
			if ok {
				err = meta.InitDBTable(vals)
			}
			break
		case "column":
			vals, ok := value.(map[string]any)
			if ok {
				err = meta.InitColumn(vals)
			}
			break
		case "query":
			vals, ok := value.(map[string]any)
			if ok {
				err = meta.InitQuery(vals)
			}
			break
		case "lang":
			vals, ok := value.(map[string]any)
			if ok {
				err = meta.InitLang(vals)
			}
			break
		case "param":
			vals, ok := value.(map[string]any)
			if ok {
				err = meta.InitParam(vals)
			}
			break
		case "convert":
			if vals, ok := value.(map[string]any); ok {
				err = meta.InitConvert(vals)
			}
		}

		if err != nil {
			return nil, err
		}
	}
	return meta, err
}
