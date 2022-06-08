package mapper

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"testing"
	"time"
)

const (
	DateFormat = "2006-01-02"
	TimeFormat = "2006-01-02 15:04:05"
)

type (
	S1 struct{ Value *time.Time }
	T1 struct{ Value *types.JSONTime }

	S2 struct{ Value time.Time }
	T2 struct{ Value *types.JSONTime }

	S3 struct{ Value types.JSONTime }
	T3 struct{ Value *time.Time }

	S4 struct{ Value types.JSONTime }
	T4 struct{ Value time.Time }
)

func Test_MapperTime1(t *testing.T) {
	SetEnabledAutoTypeConvert(true)
	valueTime := time.Now()

	source := S1{Value: &valueTime}
	var target T1
	if err := AutoMapper(&source, &target); err != nil {
		t.Error(err)
	}
}

func Test_MapperTime2(t *testing.T) {
	SetEnabledAutoTypeConvert(true)
	valueTime := time.Now()

	source := S2{Value: valueTime}
	var target T2
	if err := AutoMapper(&source, &target); err != nil {
		t.Error(err)
	}
	println(target.Value.String())
}

func Test_MapperTime3(t *testing.T) {
	SetEnabledAutoTypeConvert(true)
	valueTime := types.JSONTime(time.Now())

	source := S3{Value: valueTime}
	var target T3
	if err := AutoMapper(&source, &target); err != nil {
		t.Error(err)
	}
	println(target.Value.String())
}

func Test_MapperTime4(t *testing.T) {
	SetEnabledAutoTypeConvert(true)
	valueTime := types.JSONTime(time.Now())

	source := S4{Value: valueTime}
	var target T4
	if err := AutoMapper(&source, &target); err != nil {
		t.Error(err)
	}
	println(target.Value.String())
}
