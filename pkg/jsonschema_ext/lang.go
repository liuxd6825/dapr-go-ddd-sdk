package jsonschema_ext

import "fmt"

type Lang map[string]string

func (l Lang) init(values map[string]any) error {
	for k, v := range values {
		l[k] = fmt.Sprintf("%s", v)
	}
	return nil
}
