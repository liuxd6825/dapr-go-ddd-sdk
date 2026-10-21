package service

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"
	"github.com/stretchr/testify/assert"
)

func TestRecordService_CreateSparkTask(t *testing.T) {
	opts := CreateSparkTaskOptions{
		SparkOption: livyserver.SparkOption{
			JavaApplicationClassName: "org.example.Application",
			JavaJarFile:              "s3a://spark/app/spark-import-doris-hugegraph-4.0-SNAPSHOT.jar",
			S3: livyserver.S3Option{
				Endpoint: "192.168.120.200:8090",
				User:     "minioadmin",
				Password: "minioadmin",
			},
		},
		Args: CreateSparkTaskArgs{
			TaskId:     "test",
			CaseId:     "1001",
			S3DataFile: "s3a://import-data/record/record.parquet",
			Doris: livyserver.MySQLOption{
				Host:     "192.168.120.200",
				Port:     "8030",
				User:     "root",
				Password: "12345678",
				DB:       "master",
				Table:    "master_record",
			},
		},
	}
	err := recordService.CreateSparkTask(ctx, opts)
	assert.Nil(t, err)

}
