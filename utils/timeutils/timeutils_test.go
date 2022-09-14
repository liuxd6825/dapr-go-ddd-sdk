package timeutils

import (
	"errors"
	"testing"
	"time"
)

func TestEqual(t *testing.T) {
	t1 := time.Now()
	t2 := t1.AddDate(0, 0, 1)

	if ok := Equal(nil, nil); !ok {
		t.Error(errors.New("equal(nil, nil) error"))
	}

	if ok := Equal(t1, t2); ok {
		t.Error(errors.New("equal(t1, t2+1day) error"))
	}

	if ok := Equal(&t1, &t2); ok {
		t.Error(errors.New("equal(&t1, &t2+1day) error"))
	}

	t2 = t1
	if ok := Equal(t1, t2); !ok {
		t.Error(errors.New("equal(t1, t2) error"))
	}

	if ok := Equal(&t1, &t2); !ok {
		t.Error(errors.New("equal(&t1, &t2) error"))
	}

	if ok := Equal(t1, nil); ok {
		t.Error(errors.New("equal(t1, nil) error"))
	}

	if ok := Equal(nil, t2); ok {
		t.Error(errors.New("equal(nil, t2) error"))
	}

	if ok := Equal(t1, "nil"); ok {
		t.Error(errors.New(`equal(t1, "nil") error1)`))
	}
}
