package dao

import (
	"strings"
	"testing"
)

func TestRecordNebulaDao_BuildInsertNGQLs(t *testing.T) {
	batch := []Row{
		{
			TenantId: "t001",
			CaseId:   "case1",
			Acct:     "acct1",
			OppAcct:  "acct2",
			Date:     1722470400000,
			Name:     "张三",
			OppName:  "李四",
			Amount:   1000.0,
			Payout:   0.0,
			Income:   1000.0,
			Ccy:      "CNY",
		},
	}
	p := NebulaImportParams{
		TenantId:  "t001",
		CaseId:    "case1",
		Bucket:    "bucket1",
		Key:       "key1",
		BatchSize: 500,
	}
	ngqls := BuildInsertNGQLs(p, batch)

	if len(ngqls) == 0 {
		t.Fatal("expected non-empty nGQL list")
	}
	joined := strings.Join(ngqls, "\n")
	checks := []string{
		"INSERT VERTEX human",
		"INSERT VERTEX account",
		"INSERT EDGE owner",
		"INSERT EDGE record",
		"ON DUPLICATE KEY UPDATE",
		"txn_count = record.txn_count + 1",
		`"acct1"->"acct2"`,
	}
	for _, want := range checks {
		if !strings.Contains(joined, want) {
			t.Errorf("nGQL missing fragment: %q\n--- generated ---\n%s", want, joined)
		}
	}
}