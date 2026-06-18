package doris

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/miniofs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// These mirror the production values from cmd/hserver/config/dapr/dev_lxd.
// The test is integration-only and self-skips when the server is unreachable
// so it can run in offline environments without failing the build.
const (
	testMinIOEndpoint  = "192.168.120.224:9000"
	testMinIOAccessKey = "minioadmin"
	testMinIOSecretKey = "minioadmin"
	testMinIOUseSSL    = false
	testMinIOBucket    = "import-data"
)

// newTestMinIOClient returns a minio.Client using the same credentials as
// production, and tries both useSSL settings so the test works whether the
// server is configured for plain HTTP or HTTPS.
func newTestMinIOClient(t *testing.T) *minio.Client {
	t.Helper()
	tryConnect := func(secure bool) (*minio.Client, error) {
		c, err := minio.New(testMinIOEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(testMinIOAccessKey, testMinIOSecretKey, ""),
			Secure: secure,
		})
		if err != nil {
			return nil, err
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_, err = c.BucketExists(ctx, testMinIOBucket)
		return c, err
	}
	if c, err := tryConnect(testMinIOUseSSL); err == nil {
		return c
	} else if c, err2 := tryConnect(!testMinIOUseSSL); err2 == nil {
		t.Logf("configured useSSL=%t did not work (%v); falling back to useSSL=%t",
			testMinIOUseSSL, err, !testMinIOUseSSL)
		return c
	} else {
		t.Skipf("minio %s not reachable (ssl=%t: %v; ssl=%t: %v)",
			testMinIOEndpoint, testMinIOUseSSL, err, !testMinIOUseSSL, err2)
		return nil
	}
}

// TestParquetFile_MinIO is a regression test for the bug where ParquetFile.Close()
// forgot to close the underlying afero.File. With a streaming Fs such as
// miniofs this caused the io.Pipe to stay open forever, the PutObject
// goroutine to hang, and the object to never appear in MinIO. The
// subsequent loader.ImportFile call would then fail with NoSuchKey.
//
// The test follows the same code path as Import2Master: it creates a
// ParquetFile backed by a miniofs, writes a handful of rows, calls
// WriteStop, and then re-opens the object to confirm the upload actually
// completed and the bytes landed in the bucket.
func TestParquetFile_MinIO(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping minio integration test in -short mode")
	}
	client := newTestMinIOClient(t)
	fs, err := miniofs.NewFs(miniofs.Config{Bucket: testMinIOBucket}, client)
	if err != nil {
		t.Fatalf("miniofs.NewFs: %v", err)
	}

	path := fmt.Sprintf("/record/test_parquet_%d.parquet", time.Now().UnixNano())
	t.Cleanup(func() { _ = fs.RemoveAll("/record") })

	// 1. Drive the production code path.
	if err := fs.Mkdir("/record", 0o755); err != nil {
		t.Fatalf("Mkdir /record: %v", err)
	}
	pf, err := NewParquetFile[*Record](fs, path)
	if err != nil {
		t.Fatalf("NewParquetFile: %v", err)
	}

	now := time.Now()
	const rowCount = 10
	for i := 0; i < rowCount; i++ {
		if err := pf.Write(&Record{
			Id:          idutils.NewId(),
			TenantId:    "tenant-minio",
			CaseId:      "case-minio",
			CreatorId:   "c",
			CreatorName: "creator",
			UpdaterId:   "u",
			UpdaterName: "updater",
			Remark:      fmt.Sprintf("row-%d", i),
			DocId:       "doc",
			TaskId:      "task",
			Name:        fmt.Sprintf("name-%d", i),
			Acct:        fmt.Sprintf("acct-%d", i),
			AcctType:    "type",
			BankName:    "中行",
			Balance:     randomutils.PFloat64(),
			Payout:      randomutils.PFloat64(),
			Income:      randomutils.PFloat64(),
			Amount:      randomutils.PFloat64(),
			OppAcctType: "opp-type",
			OppAcct:     "opp-acct",
			OppName:     "opp-name",
			OppBankName: "opp-bank",
			Cash:        i%2 == 0,
			Serial:      "serial",
			MasterType:  "mt",
			MasterName:  "mn",
			MasterId:    "mi",
			GraphId:     "g",
			Month:       int32(6),
			Year:        int32(2024),
			Day:         int32(16),
			Io:          int32(1),
			Date:        &now,
			CreatedTime: &now,
			UpdatedTime: nil,
			FileId:      "file",
			SheetId:     "sheet",
			Place:       "place",
			Summary:     "summary",
			Notes:       "notes",
			Ccy:         "CNY",
			RowNum:      int32(i + 1),
		}); err != nil {
			t.Fatalf("pf.Write row %d: %v", i, err)
		}
	}

	// 2. The crucial step. Before the fix this returned nil but left the
	// streaming PutObject goroutine hanging on an open pipe, so the object
	// never reached MinIO.
	if err := pf.WriteStop(); err != nil {
		t.Fatalf("WriteStop: %v", err)
	}

	// 3. Stat must now report a real object (not a NoSuchKey error).
	info, err := fs.Stat(path)
	if err != nil {
		t.Fatalf("Stat after WriteStop: %v (object was never uploaded)", err)
	}
	if info.IsDir() {
		t.Fatalf("Stat reported a directory for %s", path)
	}
	if info.Size() == 0 {
		t.Fatalf("uploaded parquet file is 0 bytes")
	}
	t.Logf("✓ uploaded %d bytes", info.Size())

	// 4. Re-open and validate the bytes look like a real Parquet file
	// (PAR1 header + PAR1 footer per the Parquet spec).
	f, err := fs.Open(path)
	if err != nil {
		t.Fatalf("Open after WriteStop: %v", err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(data) < 8 {
		t.Fatalf("parquet file too short: %d bytes", len(data))
	}
	if !bytes.HasPrefix(data, []byte("PAR1")) {
		t.Errorf("missing PAR1 header, got %q", data[:4])
	}
	if !bytes.HasSuffix(data, []byte("PAR1")) {
		t.Errorf("missing PAR1 footer, got %q", data[len(data)-4:])
	}
}

// TestParquetFile_MinIO_CloseIsIdempotent ensures that calling Close() twice
// does not panic and does not double-report an error, since production code
// may invoke it defensively.
func TestParquetFile_MinIO_CloseIsIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping minio integration test in -short mode")
	}
	client := newTestMinIOClient(t)
	fs, err := miniofs.NewFs(miniofs.Config{Bucket: testMinIOBucket}, client)
	if err != nil {
		t.Fatalf("miniofs.NewFs: %v", err)
	}
	t.Cleanup(func() { _ = fs.RemoveAll("/record") })

	path := fmt.Sprintf("/record/test_idempotent_%d.parquet", time.Now().UnixNano())
	if err := fs.Mkdir("/record", 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	pf, err := NewParquetFile[*Record](fs, path)
	if err != nil {
		t.Fatalf("NewParquetFile: %v", err)
	}
	if err := pf.WriteStop(); err != nil {
		t.Fatalf("first WriteStop: %v", err)
	}
	// Second call must be a no-op.
	if err := pf.WriteStop(); err != nil {
		t.Errorf("second WriteStop returned an error: %v", err)
	}
}
