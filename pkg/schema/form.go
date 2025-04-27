package schema

type Form map[string]any

func (c Form) init(vals map[string]any) error {
	for k, v := range vals {
		c[k] = v
	}
	return nil
}

func NewForm() *Form {
	return &Form{}
}
