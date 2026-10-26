package yamlutils

import (
	"gopkg.in/yaml.v3"
)

func Marshal(data any) ([]byte, error) {
	yamlBytes, err := yaml.Marshal(data)
	return yamlBytes, err
}

func Unmarshal(yamlBytes []byte, data any) error {
	err := yaml.Unmarshal(yamlBytes, data)
	return err
}

func Copy(src, dst interface{}) error {
	bytes, err := Marshal(src)
	if err != nil {
		return err
	}
	return Unmarshal(bytes, dst)
}
