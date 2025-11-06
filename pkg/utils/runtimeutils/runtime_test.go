package runtimeutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoId(t *testing.T) {
	goid, err := GoId()
	assert.NoError(t, err)
	t.Log("goId", goid)
}

func TestGetFuncName(t *testing.T) {
	funcName := GetFuncName(0)
	t.Log("TestGetFuncName", funcName)
}

func TestGetPackageName(t *testing.T) {
	pkgName := GetPackageName(0)
	t.Log("TestGetPackageName", pkgName)
}
