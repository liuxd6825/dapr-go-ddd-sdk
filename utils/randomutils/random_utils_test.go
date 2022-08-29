package randomutils

import (
	"testing"
)

func Test_String(t *testing.T) {
	t.Logf("String() : %v", String(10))
	t.Logf("StringNumber() : %v", StringNumber(10))
	t.Logf("StringLower() : %v", StringLower(10))
	t.Logf("StringUpper() : %v", StringUpper(10))
	t.Logf("ChinaName() : %v", ChinaName())
	t.Logf("Email() : %v", Email())
	t.Logf("UUID() : %v", UUID())
	t.Logf("DateString() : %v ", DateString())
	for i := 0; i < 100; i++ {
		t.Logf("TimeString() : %v ", TimeString())
	}

	for i := 0; i < 100; i++ {
		t.Logf("Time() : %v", Time())
	}

	t.Logf("IpAddr() : %v ", IpAddr())

	t.Logf("Int() : %v ", Int())
	t.Logf("Int64(1) : %v ", Int64())
	t.Logf("Float32() : %v ", Float32())
	t.Logf("Float64() : %v ", Float64())
}
