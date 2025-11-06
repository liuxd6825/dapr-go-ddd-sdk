package schema

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store"
	"github.com/liuxd6825/jsonschema/v6"
)

func PagingQueryConverter(data any, sch *jsonschema.Schema, meta *MetaExtension) (any, error) {
	mapData := data.(map[string]interface{})
	builder := store.NewFindPagingQueryBuilder().SetMapToQuery(mapData)
	return builder.Build(), nil
}
