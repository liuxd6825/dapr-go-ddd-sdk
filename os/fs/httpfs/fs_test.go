package httpfs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/os/fs/fsm"
	"testing"
)

func TestConnect(t *testing.T) {
	service := fsm.NewManager()

	if err := service.Open(); err != nil {
		t.Error(err)
	}
	if err := service.Open(); err != nil {
		t.Error(err)
	}
}
