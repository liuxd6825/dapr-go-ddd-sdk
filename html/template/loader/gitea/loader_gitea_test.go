package gitea

import (
	"testing"
)

func TestConnect(t *testing.T) {
	service := NewManage("schema", "draw-web", "main")
	if err := service.Connect("http://localhost:3000", "liuxd", "liuxd"); err != nil {
		t.Error(err)
	}
	if _, err := service.GetFile("schema/master/human/human.schema.json"); err != nil {
		t.Error(err)
	}
}
