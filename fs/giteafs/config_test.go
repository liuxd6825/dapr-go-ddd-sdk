package giteafs

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	meta := map[string]any{
		"name":     "/path/to/repo",
		"url":      "https://localhost:3000",
		"user":     "liuxd",
		"password": "liuxd",
		"email":    "admin@localhost",
		"repo":     "test",
		"branch":   "main",
	}

	config, err := NewConfig(meta)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(config)
}
