package fs

import (
	"testing"
)

func Test_NewManagerWithConfigs(t *testing.T) {
	medata := []map[string]any{
		{"id": "gitea", "type": "gitea", "user": "liuxd", "password": "liuxd", "url": "http://localhost:3000", "repo": "test", "branch": "main", "email": "admin@liuxd.com"},
		{"id": "local", "type": "local", "path": "/Users/lxd/Projects/liuxd6825/dapr/dapr-go-ddd-sdk/fs/test"},
	}
	manager, err := NewManagerWithConfigs(medata)
	if err != nil {
		t.Error(err)
	}
	fileName := "local:/text.txt"
	if err := manager.WriteFile(fileName, "", []byte("123"), WriteModelAllWriteRead); err != nil {
		t.Error(err)
	}

	if bytes, err := manager.ReadFile(fileName, ""); err != nil {
		t.Error(err)
	} else {
		t.Log(string(bytes))
	}

}
