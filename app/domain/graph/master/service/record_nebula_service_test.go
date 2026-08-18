package service

import (
	"context"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

func TestRecordNebulaService_EmptyPath(t *testing.T) {
	svc := &RecordNebulaService{dao: nil, env: &env.Env{}}
	ctx := context.Background()
	if _, err := svc.ImportFromS3(ctx, "", "case1", ImportOptions{}); err == nil {
		t.Fatal("expected error for empty s3Path")
	}
	if _, err := svc.ImportFromS3(ctx, "s3://x", "", ImportOptions{}); err == nil {
		t.Fatal("expected error for empty caseId")
	}
}

func TestRecordNebulaService_ResolveS3Path(t *testing.T) {
	e := &env.Env{Minio: map[string]*env.Minio{
		"default": {
			AccessKey: "ak", SecretKey: "sk", Endpoint: "host:9000",
			Buckets: map[string]string{"importData": "import-data"},
		},
	}}
	svc := &RecordNebulaService{env: e}

	cases := []struct {
		in      string
		wantBk  string
		wantKey string
		wantErr bool
	}{
		{"s3://ak:sk@host/bucket/key1.parquet", "bucket", "key1.parquet", false},
		{"s3://bucket/key2.parquet", "bucket", "key2.parquet", false},
		{"key3.parquet", "import-data", "key3.parquet", false},
		{"", "", "", true},
	}
	for _, c := range cases {
		bk, k, err := svc.resolveS3Path(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("input %q: expected error", c.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("input %q: %v", c.in, err)
		}
		if bk != c.wantBk {
			t.Errorf("bucket: got %q want %q", bk, c.wantBk)
		}
		if k != c.wantKey {
			t.Errorf("key: got %q want %q", k, c.wantKey)
		}
	}
}
