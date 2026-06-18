package service

import (
	"context"
	"strings"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

const (
	testS3Host     = "192.168.120.224:9000"
	testS3Access   = "minioadmin"
	testS3Secret   = "minioadmin"
	testBucket     = "import-data"
	testTenantId   = "test"
	testCaseId     = "PnRAAeb4liYYMyLfpsA7e9bu"
	testParquetKey = "record/record_hUe2yxR81YNB89AY3sWRIM4U.parquet"
)

// TestRecordParquetService_BuildCypher 离线断言: 生成的 Cypher 包含关键片段
func TestRecordParquetService_BuildCypher(t *testing.T) {
	p := dao.ImportParams{
		TenantId:    testTenantId,
		CaseId:      testCaseId,
		S3Path:      "s3://x",
		BatchSize:   500,
		Parallel:    false,
		Concurrency: 1,
		Retries:     1,
	}
	cypher := dao.BuildImportCypher(p)

	checks := []string{
		"apoc.periodic.iterate",
		"apoc.load.parquet",
		"date(datetime(row.date))",
		"MERGE (a1)-[r:record",
		"'rec:' + tenantId + ':' + caseId",
		"ON CREATE SET",
		"ON MATCH SET",
		"r.txn_count    = coalesce(r.txn_count, 0) + 1",
		"parallel: false",
		"concurrency: 1",
		"retries: $retries",
	}
	for _, want := range checks {
		if !strings.Contains(cypher, want) {
			t.Errorf("Cypher missing fragment: %q", want)
		}
	}
}

// TestRecordParquetService_BuildS3URL 验证 s3 URL 拼接
func TestRecordParquetService_BuildS3URL(t *testing.T) {
	m := &env.Minio{
		AccessKey: testS3Access,
		SecretKey: testS3Secret,
		Endpoint:  testS3Host,
	}
	url := dao.BuildS3URL(m, testBucket, testParquetKey)
	expected := "s3://minioadmin:minioadmin@192.168.120.224:9000/import-data/record/record_hUe2yxR81YNB89AY3sWRIM4U.parquet"
	if url != expected {
		t.Errorf("BuildS3URL wrong:\n  got:  %s\n  want: %s", url, expected)
	}
}

// TestRecordParquetService_EmptyPath 错误路径 (不依赖 Neo4j/MinIO)
func TestRecordParquetService_EmptyPath(t *testing.T) {
	svc := &RecordParquetService{dao: nil, env: nil}
	ctx := context.Background()
	_, err := svc.ImportFromS3(ctx, "", "case1", ImportOptions{})
	if err == nil {
		t.Fatal("expected error for empty s3Path")
	}
	_, err = svc.ImportFromS3(ctx, "s3://x", "", ImportOptions{})
	if err == nil {
		t.Fatal("expected error for empty caseId")
	}
}

// TestRecordParquetService_ImportFromS3 集成测试: 需 MinIO + Neo4j 可达
func TestRecordParquetService_ImportFromS3(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in -short mode")
	}

	envVal := xtest2.InitEnv_Neo4j(xtest2.Neo4jRemoveOption)
	ctx := xtest2.NewContext()
	svc := NewRecordParquetService(envVal)

	s3Path := "s3://" + testS3Access + ":" + testS3Secret + "@" +
		testS3Host + "/" + testBucket + "/" + testParquetKey

	// 第一次导入
	res, err := svc.ImportFromS3(ctx, s3Path, testCaseId, ImportOptions{
		BatchSize: 500,
	})
	if err != nil {
		t.Fatalf("first ImportFromS3 failed: %v", err)
	}
	if res.Total != 738 {
		t.Errorf("expected total=738, got %d", res.Total)
	}
	if res.FailedOperations != 0 {
		t.Errorf("expected failed=0, got %d (errors: %v)",
			res.FailedOperations, res.ErrorMessages)
	}
	if res.Batches < 1 {
		t.Errorf("expected batches >= 1, got %d", res.Batches)
	}

	// 第二次导入 (幂等性)
	res2, err := svc.ImportFromS3(ctx, s3Path, testCaseId, ImportOptions{
		BatchSize: 500,
	})
	if err != nil {
		t.Fatalf("second ImportFromS3 failed: %v", err)
	}
	if res2.Total != 738 {
		t.Errorf("second run expected total=738, got %d", res2.Total)
	}

	// 验证: sum(r.txn_count) 应 = 2 * 738 = 1476
	sumCount := verifySumCount(t, ctx, testTenantId)
	if sumCount != 1476 {
		t.Errorf("expected sumCount=1476, got %d", sumCount)
	}
}

// verifySumCount 验证 sum(r.txn_count) 的辅助函数
func verifySumCount(t *testing.T, ctx context.Context, tenantId string) int {
	t.Helper()
	d := dao.NewRecordParquetDao()
	verifyCypher := `
		MATCH ()-[r:record {tenant_id: $tid}]->()
		RETURN sum(r.txn_count) AS s`
	result, err := d.GetStore().Write(ctx, verifyCypher, map[string]any{
		"tid": tenantId,
	})
	if err != nil {
		t.Fatalf("verify query failed: %v", err)
	}
	if result == nil {
		return 0
	}
	data := result.Data()
	if v, ok := data["s"]; ok && len(v) > 0 {
		switch x := v[0].(type) {
		case int64:
			return int(x)
		case int:
			return x
		case float64:
			return int(x)
		}
	}
	return 0
}
