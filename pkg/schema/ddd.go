package schema

import "github.com/liuxd6825/dapr-go-ddd-sdk/utils/convert"

type DDD struct {
	AggField   string `json:"aggField"`
	AggType    string `json:"aggType"`
	IsPubEvent bool   `json:"isPubEvent"`
}

func (p *DDD) init(values map[string]any) error {
	for k, v := range values {
		switch k {
		case "aggField":
			p.AggField = v.(string)
		case "aggType":
			p.AggType = v.(string)
		case "isPubEvent":
			if val, err := convert.ConvertBool(v); err == nil {
				p.IsPubEvent = val
			} else {
				return err
			}
		}
	}
	return nil
}
