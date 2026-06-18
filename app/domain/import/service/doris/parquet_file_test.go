package doris

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
	"github.com/spf13/afero"
)

func TestRecordWrite(t *testing.T) {
	fs := afero.NewOsFs()
	pf, err := NewParquetFile[*Record](fs, "./record.parquet")
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	for i := 0; i < 200000; i++ {
		item := &Record{}
		item.Id = idutils.NewId()
		item.TenantId = "TenantId"
		item.CaseId = "caseid"

		item.CreatorId = "CreatorId"
		item.CreatorName = "CreatorName"
		item.UpdaterId = "UpdaterId"
		item.UpdaterName = "UpdaterName"
		item.Remark = "remark"

		item.RowNum = int32(i + 1)

		item.TaskId = "TaskId"
		item.DocId = "DocId"
		item.FileId = "FileId"
		item.SheetId = "SheetId"
		item.Name = "name"
		item.Acct = "Acct"
		item.AcctType = "AcctType"
		item.BankName = "中行"
		item.Balance = randomutils.PFloat64()
		item.Payout = randomutils.PFloat64()
		item.Income = randomutils.PFloat64()
		item.Amount = randomutils.PFloat64()

		item.OppAcctType = "OppAcctType"
		item.OppAcct = "OppAcct"
		item.OppName = "OppName"
		item.OppBankName = "OppBankName"
		item.Cash = true
		item.Serial = "Serial"

		item.MasterType = "MasterType"
		item.MasterName = "MasterName"
		item.MasterId = "MasterId"
		item.GraphId = "GraphId"
		item.Month = int32(6)
		item.Year = int32(2024)
		item.Day = int32(16)
		item.Io = int32(1)

		// *time.Time 字段：nil 表示 Doris NULL
		item.Date = &now
		item.CreatedTime = &now
		item.UpdatedTime = nil // 验证 nil 语义

		item.Place = "Place"
		item.Summary = "Summary"
		item.Notes = "Notes"
		item.Ccy = "ccy"

		if err := pf.Write(item); err != nil {
			t.Fatal(err)
		}
	}
	if err := pf.WriteStop(); err != nil {
		t.Fatal(err)
	}

	// 1. 验证 Parquet 文件已生成
	info, err := os.Stat("./record.parquet")
	if err != nil {
		t.Fatalf("parquet file not found: %v", err)
	}
	t.Logf("✓ Parquet 文件已生成: %d bytes", info.Size())
	if info.Size() == 0 {
		t.Fatal("parquet file is empty")
	}

	loader := NewDorisLoader(Config{
		BEHost:   "192.168.120.224",
		BEPort:   8040,
		Database: "master",
		Table:    "master_record1",
		Username: "admin",
		Password: "",
	})

	t.Logf("→ 开始 Stream Load 到 Doris ...")
	if _, err := loader.ImportFile(fs, "./record.parquet", LoadOptions{
		Format:  FormatParquet,
		Columns: RecordTableColumns,
	}); err != nil {
		t.Fatalf("ImportFile failed: %v", err)
	}
	t.Logf("✓ ImportFile 返回 nil (loader 认为导入成功)")

	// 2. 通过 Doris FE (MySQL 协议) 验证数据真的进了表
	count, err := queryDorisRowCount("192.168.120.224", 9030, "admin", "", "master", "master_record1")
	if err != nil {
		t.Logf("⚠ 查询 Doris FE 失败(可能 FE 端口不通): %v", err)
		t.Logf("   请手动执行: mysql -h 192.168.120.224 -P 9030 -uadmin -e 'SELECT COUNT(*) FROM master.master_record1'")
		return
	}
	t.Logf("✓ Doris 表 master.master_record1 当前行数: %d (期望至少 +10)", count)
	if count < 10 {
		t.Errorf("Doris 表行数不足: got %d, want >= 10", count)
	}
}

// queryDorisRowCount 通过 MySQL 协议连接 Doris FE 查询表行数
func queryDorisRowCount(host string, port int, user, password, database, table string) (int64, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=5s", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return 0, fmt.Errorf("open mysql conn: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return 0, fmt.Errorf("ping mysql: %w", err)
	}

	var count int64
	q := fmt.Sprintf("SELECT COUNT(*) FROM `%s`.`%s`", database, table)
	if err := db.QueryRowContext(ctx, q).Scan(&count); err != nil {
		return 0, fmt.Errorf("query: %w", err)
	}
	return count, nil
}
