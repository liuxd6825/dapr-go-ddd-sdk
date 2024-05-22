package httpfs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"testing"
)

func TestConnect(t *testing.T) {
	service := fs.NewManager()

	if err := service.Open(); err != nil {
		t.Error(err)
	}
	if err := service.Open(); err != nil {
		t.Error(err)
	}
}
