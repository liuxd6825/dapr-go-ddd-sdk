package idutils

import (
	"testing"
)

func TestNewId(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := NewId()
		t.Log(id, " = ", len(id))
	}
}
