package idutils

import (
	"testing"
)

func TestNewId(t *testing.T) {
	id := NewId()
	t.Log(id)
}
