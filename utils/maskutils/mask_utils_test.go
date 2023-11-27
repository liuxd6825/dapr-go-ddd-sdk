package maskutils

import "testing"

func TestMask_Contain(t *testing.T) {
	values := []string{
		"AccountCode",
		"BankName",
		"EndTime",
		"FileId",
	}
	mask := NewUpdateMask(values)
	if mask.Contain("AccountCode") {
		t.Log("AccountCode")
	}
	if mask.Contain("BankName") {
		t.Log("BankName")
	}
	if mask.Contain("EndTime") {
		t.Log("EndTime")
	}
	if mask.Contain("FileId") {
		t.Log("FileId")
	}
}
