package schema_utils

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/fsm"
	"net/url"
)

type SchemaLoader struct {
	fsm *fsm.Manager
}

func NewSchemaLoader(fsm *fsm.Manager) *SchemaLoader {
	return &SchemaLoader{fsm: fsm}
}

func (l *SchemaLoader) Load(refUrl string) (any, error) {
	// 解析路径
	refPath, err := url.Parse(refUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ref path: %w", err)
	}

	// 处理相对路径
	data, err := l.fsm.ReadFile(refPath.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	// 解析 JSON 数据
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return jsonData, nil
}
