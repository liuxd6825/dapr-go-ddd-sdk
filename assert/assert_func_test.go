package assert

import "testing"

func TestNil(t *testing.T) {
	if err := Nil(nil); err != nil {
		t.Error(err)
	}

	if err := NotNil(nil); err != nil {
		t.Error(err)
	}

	a := &Object{Value: "1"}
	b := &Object{Value: "1"}
	if err := Equal(a, b); err != nil {
		t.Error(err)
	}

	if err := NotEqual(a, b); err != nil {
		t.Error(err)
	}

	s1 := "111"
	s2 := "111"
	if err := Equal(s1, s1, NewOptions("Equal(s1, s1)")); err != nil {
		t.Error(err)
	}

	if err := NotEqual(s1, s2, NewOptions("NotEqual(s1, s2))")); err != nil {
		t.Error(err)
	}
}

type Object struct {
	Value string
}
