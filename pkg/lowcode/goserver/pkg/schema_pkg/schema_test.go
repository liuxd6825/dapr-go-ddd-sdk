package schema_pkg

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/lowcode/goserver/pkg/fs_pkg"
	"github.com/open-policy-agent/opa/sdk/test"
)

func TestSchemaPkg_LoadFile(t *testing.T) {
	server, err := test.NewServer("/testfile", func(server *test.TestServer) error {
		fsPkg, err := fs_pkg.NewFsPkg(server.EnvConfig())
		server.FsPkg = fsPkg
		return err
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	schema := LoadFile("/human.json", "")

	data := map[string]any{
		"id":       "0",
		"name":     "John",
		"lastName": "Doe",
		"email":    "john.doe@example.com",
		"phone":    "123-456-7890",
	}

	// 验证数据
	err = schema.Validate(data)
	if err != nil {
		t.Errorf("Validation failed: %v\n", err)
	} else {
		t.Log("Validation successful!")
	}

}
