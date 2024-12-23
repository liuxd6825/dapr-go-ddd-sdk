package schema_pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/test"
	"testing"
)

func TestSchemaPkg_LoadFile(t *testing.T) {
	server, err := test.NewServer("")
	if err != nil {
		t.Fatal(err)
		return
	}
	schemaPkg := New(server)
	schema := schemaPkg.LoadFile("/testfile/human.json", "")

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
