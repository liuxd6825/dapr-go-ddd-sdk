package timeutils

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/setting"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStrToDateTime(t *testing.T) {
	dstr := "2021-08-18 00:00:00"
	tstr := "23595969"
	dt := time.Date(2021, 8, 18, 23, 59, 59, 690000000, time.Local)
	v, err := ToDateTime(dstr, tstr)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, v, dt)
}

func TestNow(t *testing.T) {
	setting.SetLocalTimeZone()
	now := Now().UTC()
	t.Logf("now=%v", now)
}

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
