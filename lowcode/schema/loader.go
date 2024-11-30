package schema

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"net/url"
)

// URLLoader knows how to load json from given url.
type URLLoader interface {
	// Load loads json from given absolute url.
	Load(url string) (any, error)
}

type JSONLoader struct {
	fsm *fs.Manager
}

func NewJSONLoader(fsm *fs.Manager) *JSONLoader {
	return &JSONLoader{fsm: fsm}
}

func (cl *JSONLoader) Load(refUrl string) (any, error) {
	// 解析路径
	refPath, err := url.Parse(refUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ref path: %w", err)
	}

	// 处理相对路径
	data, err := cl.fsm.ReadFile(refPath.Path)
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
