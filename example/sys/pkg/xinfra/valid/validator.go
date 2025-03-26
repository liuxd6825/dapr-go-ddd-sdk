package valid

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
)

func GenerateSchema(s interface{}) ([]byte, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties:  false, // 禁止额外字段
		RequiredFromJSONSchemaTags: true,  // 从 JSON Schema 标签解析 required
	}
	schema := reflector.Reflect(s)
	return json.MarshalIndent(schema, "", "  ")
}
