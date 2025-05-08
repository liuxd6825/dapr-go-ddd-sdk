package definition

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"os"
	"path/filepath"
	"testing"
)

func TestDefinition_NewDefinition(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
		return
	}

	srcPath := filepath.Join(path, "testfile")

	fs, err := localfs.NewFs(localfs.Config{Name: "file", Path: srcPath})
	if err != nil {
		t.Fatal(err)
		return
	}

	d, err := NewDefinition(fs, "/definition/params")
	if err != nil {
		t.Fatal(err)
		return
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(data))
}
