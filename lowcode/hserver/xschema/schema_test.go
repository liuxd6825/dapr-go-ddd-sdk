package xschema

import (
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/localfs"
	"github.com/liuxd6825/jsonschema/v6"
	"github.com/spf13/afero"
	"net/url"
	"os"
	"testing"
)

func TestSchema_GetJsonSchema(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
		return
	}
	path = path + "/testfile"

	fs, err := localfs.NewFs(localfs.Config{Name: "file", Path: path})
	if err != nil {
		t.Fatal(err)
		return
	}

	fileData, err := afero.ReadFile(fs, "/human.json")
	if err != nil {
		t.Fatal(err)
	}

	scm := NewSchema()
	err = json.Unmarshal(fileData, scm)
	if err != nil {
		t.Fatal(err)
	}

	scm.AllOf = []*Schema{
		{Ref: "/base.json"},
	}
	loader := newUrlLoader(fs)
	schema := scm.Init("human.json", loader)

	props := schema.GetAllProperties()
	for key, prop := range props {
		t.Log(key, prop.Title)
	}

	t.Log("//  schema.GetFields() ")

	dateFields := schema.GetFields(jsonschema.JsonType_DateType, jsonschema.JsonType_DateTimeType)
	printMap("", dateFields, t)
}

func printMap(parentKey string, m map[string]any, t *testing.T) {
	for key, field := range m {
		if prop, ok := field.(*jsonschema.Schema); ok {
			t.Log(parentKey+"."+key, prop.Title)
		} else if fmap, ok := field.(map[string]any); ok {
			printMap(parentKey+"."+key, fmap, t)
		}
	}
}

type urlLoader struct {
	fs afero.Fs
}

func newUrlLoader(fs afero.Fs) *urlLoader {
	return &urlLoader{fs: fs}
}

func (l *urlLoader) Load(refUrl string) (any, error) {
	// 解析路径
	refPath, err := url.Parse(refUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ref path: %w", err)
	}
	data, err := afero.ReadFile(l.fs, refPath.Path)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 数据
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return jsonData, nil

}
