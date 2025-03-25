package jsonschema_ext

type Control map[string]any

func (c Control) init(vals map[string]any) error {
	for k, v := range vals {
		c[k] = v
	}
	return nil
}

/*
type ControlType string
type Control struct {
	Type       ControlType    `json:"type,omitempty"`
	Properties map[string]any `json:"properties,omitempty"`
}

func (c *Control) init(values map[string]any) error {
	for k, v := range values {
		switch k {
		case "type":
			c.Type = ControlType(v.(string))
		case "properties":
			if val, ok := v.(map[string]any); ok {
				c.Properties = val
			}
		}
	}
	return nil
}
*/
