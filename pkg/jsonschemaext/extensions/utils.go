package extensions

import "github.com/liuxd6825/jsonschema/v6"

func GetTableName(sch *jsonschema.Schema) string {
	tableName := sch.Name()
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			if meta.DBTable != nil && meta.DBTable.Name != "" {
				tableName = meta.DBTable.Name
				break
			}
		}
	}
	return tableName
}

func GetFieldName(sch *jsonschema.Schema) string {
	fieldName := sch.Name()
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			if meta.DBField != nil && meta.DBField.Name != "" {
				fieldName = meta.DBField.Name
				break
			}
		}
	}
	return fieldName
}

func GetField(sch *jsonschema.Schema) *DBField {
	for _, e := range sch.Extensions {
		if meta, ok := e.(*MetaExtension); ok {
			return meta.DBField
		}
	}
	return nil
}
