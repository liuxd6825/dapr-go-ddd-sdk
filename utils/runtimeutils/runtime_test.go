package runtimeutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGoId(t *testing.T) {
	goid, err := GoId()
	assert.NoError(t, err)
	t.Log("goId", goid)
}
