package ddd_utils

import (
	"encoding/json"
	"errors"
	"fmt"
)

func NewMap() map[string]interface{} {
	return make(map[string]interface{})
}

func IsEmpty(v string, field string) error {
	if len(v) == 0 {
		return errors.New(fmt.Sprintf("%s  cannot be empty.", field))
	}
	return nil
}

func NewMapInterface(jsonText string) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonText), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ToJson(data interface{}) (string, error) {
	if data == nil {
		return "", nil
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func NewAppError(appID string, err error) error {
	msg := fmt.Sprintf("AppId is %s , %s", appID, err.Error())
	return errors.New(msg)
}
