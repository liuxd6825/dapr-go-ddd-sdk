package model

import (
	"testing"
	"time"
)

func Test_NewTranId(t *testing.T) {
	amount := 100.01
	date := time.Now()
	record := NewRecord()
	record.Acct = "acct"
	record.OppAcct = "oppAcct"
	record.Amount = &amount
	record.Income = &amount
	record.Date = &date
	tranId := NewTranId(record)
	t.Log(tranId)

	payout := 0 - amount
	record.Acct = "oppAcct"
	record.OppAcct = "acct"
	record.Amount = &amount
	record.Date = &date
	record.Payout = &payout
	tranId2 := NewTranId(record)
	t.Log(tranId2)
}
