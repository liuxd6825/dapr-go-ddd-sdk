package fs

import "testing"

func Test_NewManagerWithConfigs(t *testing.T) {
	medata := []map[string]any{
		{"id": "id1", "type": "gitea", "user": "liuxd", "password": "liuxd", "url": "http://localhost:3000", "repo": "test", "branch": "main", "email": "admin@liuxd.com"},
		{"id": "id2", "type": "local", "path": "/test2"},
	}
	manager, err := NewManagerWithConfigs(medata)
	if err != nil {
		t.Error(err)
	}
	t.Log(manager)
}
