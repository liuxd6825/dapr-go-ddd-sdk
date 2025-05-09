package schema

import "github.com/liuxd6825/jsonschema/v6"

type Column map[string]any

func (c Column) init(ctx *jsonschema.CompilerContext, vals map[string]any) error {
	for k, v := range vals {
		c[k] = v
	}
	return nil
}

func NewColumn() *Column {
	return &Column{}
}

/*
type ColumnType string

type Column struct {
	Type  ColumnType `json:"type,omitempty"`
	Width int64      `json:"width,omitempty"`
	Order int64      `json:"order,omitempty"`
	Hide  bool       `json:"hide,omitempty"`
	Flex  int64      `json:"flex,omitempty"`
}

func (c *Column) init(values map[string]any) error {
	var err error
	for k, v := range values {
		switch k {
		case "type":
			if v, ok := v.(string); ok {
				c.Type = ColumnType(v)
			}
			break
		case "width":
			if val, ok := v.(json.Number); ok {
				c.Width, err = val.Int64()
			}
			break
		case "order":
			if val, ok := v.(json.Number); ok {
				c.Order, err = val.Int64()
			}
			break
		case "hide":
			if val, ok := v.(bool); ok {
				c.Hide = val
			}
		case "flex":
			if val, ok := v.(json.Number); ok {
				c.Flex, err = val.Int64()
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
*/
