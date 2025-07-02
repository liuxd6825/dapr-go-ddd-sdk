package utils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
	"github.com/liuxd6825/jsonschema/v6"
	"strings"
)

func GetNodeType(dbSch *dbschema.DBSchema) string {
	if dbSch == nil || dbSch.JsonSchema == nil {
		return ""
	}
	meta := schema.GetMetaExtension(dbSch.JsonSchema)
	if meta == nil || meta.Attributes == nil {
		return ""
	}
	attr := getAttributes(dbSch.JsonSchema)
	if attr == nil {
		return ""
	}
	if typeVal, ok := attr["title"]; ok {
		return typeVal.(string)
	}
	if labelsVal, ok := attr["labels"]; ok {
		if list, ok := labelsVal.([]string); ok {
			return strings.Join(list, ",")
		}
	}
	return ""
}

func getAttributes(sch *jsonschema.Schema) map[string]any {
	if sch == nil {
		return nil
	}
	meta := schema.GetMetaExtension(sch)
	if meta == nil || meta.Attributes == nil {
		return nil
	}
	if val, ok := meta.Attributes["neo4j"]; ok {
		if mapVal, ok := val.(map[string]interface{}); ok {
			return mapVal
		}
	}
	return nil
}

func GetDescription(data map[string]any, dbSch *dbschema.DBSchema) string {
	if dbSch == nil {
		return ""
	}
	sb := strings.Builder{}

	for _, field := range dbSch.Fields {
		dbName := field.DBName
		if !isDescField(dbName) {
			continue
		}

		if val, ok := data[dbName]; ok {
			if val != nil {
				if field.Title != "" {
					sb.WriteString(field.Title + ":")
				} else {
					sb.WriteString(dbName + ":")
				}
				sb.WriteString(fmt.Sprintf("%v; ", val))
			}
		}
	}
	return sb.String()
}

func isDescField(fieldName string) bool {
	name := strings.ToLower(fieldName)
	if equal(name, "id", "is_deleted", "deleter_name", "updater_name", "creator_name", "relation_type") || hasSuffix(name, "_id", "_time", "_ids") {
		return false
	}
	return true
}

func equal(fieldName string, vals ...string) bool {
	for _, val := range vals {
		if val == fieldName {
			return true
		}
	}
	return false
}

func hasSuffix(fieldName string, suffix ...string) bool {
	for _, suffix := range suffix {
		if strings.HasSuffix(fieldName, suffix) {
			return true
		}
	}
	return false
}
