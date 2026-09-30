package service_test

import (
	"context"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/service"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

func TestRecordService_Import(t *testing.T) {
	s := service.NewRecordService(service.RecordServiceConfig{
		S3Endpoint:      "http://192.168.120.224:9000",
		S3User:          "minioadmin",
		S3Pwd:           "minioadmin",
		HGHost:          "192.168.120.200",
		HGPort:          "18080",
		HGGraph:         "hugegraph",
		HGUser:          "admin",
		HGPwd:           "admin",
		MongoHost:       "192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020",
		MongoReplicaset: "mongors",
		MongoDB:         "master",
		MongoUser:       "super_admin",
		MongoPwd:        "123456",
		LivyURL:         "http://192.168.120.200:8998/batches",
	})

	ctx := appctx.NewTenantContext(context.Background(), "001")
	jarFile := "s3a://spark/app/record-parquet-to-hugegraph-4.0-SNAPSHOT.jar"
	dataFile := "s3a://import-data/record/record.parquet"
	graphName := "hugegraph"
	err := s.Import(ctx, "1001", "2ndoebXy02tqJTdy9YisXp4I", dataFile, jarFile, graphName)
	if err != nil {
		t.Error(err)
	}
}
