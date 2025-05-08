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

func Test_fmtDateStr(t *testing.T) {
	data1, _ := fmtDateStr("2022年7月20日")
	assert.Equal(t, data1, "2022-07-20")

	data2, _ := fmtDateStr("20220720")
	assert.Equal(t, data2, "2022-07-20")

	data3, _ := fmtDateStr("2022/07/20")
	assert.Equal(t, data3, "2022-07-20")

	date4, _ := fmtDateStr("2022-07-20")
	assert.Equal(t, date4, "2022-07-20")

	date5, _ := fmtDateStr("2022.7.20")
	assert.Equal(t, date5, "2022-07-20")
}

func Test_fmtTimeStr(t *testing.T) {
	data1, _ := fmtTimeStr("12时10分15秒")
	assert.Equal(t, data1, "12:10:15")

	data2, _ := fmtTimeStr("12:1:1")
	assert.Equal(t, data2, "12:01:01")

	data3, _ := fmtTimeStr("12.10.10")
	assert.Equal(t, data3, "12:10:10")

	date4, _ := fmtTimeStr("12,10,10")
	assert.Equal(t, date4, "12:10:10")

	date5, _ := fmtTimeStr("12 10 10")
	assert.Equal(t, date5, "12:10:10")
}

func Test_StrToDateTime(t *testing.T) {
	v1, err := StrToDateTime("20140101 12:10:10")
	assert.NoError(t, err)
	t.Log(v1)

	v2, err := StrToDateTime("2014-01-01 12:10:10")
	assert.NoError(t, err)
	t.Log(v2)

	v3, err := StrToDateTime("2014年01月01日 12时10分10秒")
	assert.NoError(t, err)
	t.Log(v3)

	v4, err := StrToDateTime("20140101 121010")
	assert.NoError(t, err)
	t.Log(v4)

	v5, err := StrToDateTime("2014.01.01 12.10.10")
	assert.NoError(t, err)
	t.Log(v5)
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
