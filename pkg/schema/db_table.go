package schema

import "github.com/liuxd6825/jsonschema/v6"

type DBTable struct {
	Name       string         `json:"name"`
	DBKey      string         `json:"dbKey"`
	Properties map[string]any `json:"properties"`
}

func NewDBTable() *DBTable {
	return &DBTable{}
}

func (db *DBTable) init(ctx *jsonschema.CompilerContext, values map[string]any) error {
	var err error
	for k, v := range values {
		switch k {
		case "name":
			db.Name = v.(string)
		case "dbKey":
			db.DBKey = v.(string)
		case "properties":
			{
				prop, ok := v.(map[string]any)
				if ok {
					db.Properties = prop
				}
			}
		}
	}
	return err
}
