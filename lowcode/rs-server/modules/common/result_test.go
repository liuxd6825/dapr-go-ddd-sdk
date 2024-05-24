package common

import "testing"

func TestResult(t *testing.T) {
	res := NewResult[any](nil, nil)
	t.Log(res.GetFoundResults())
}
