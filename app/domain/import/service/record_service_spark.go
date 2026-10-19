package service

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"
)

type CreateSparkTaskArgs struct {
	AppID      string                 `json:"appId"`
	CaseId     string                 `json:"caseId"`
	TaskId     string                 `json:"taskId"`
	S3DataFile string                 `json:"s3DataFile"`
	Doris      livyserver.MySQLOption `json:"doris"`
}

type CreateSparkTaskOptions struct {
	livyserver.SparkOption
	Args CreateSparkTaskArgs
}

// CreateSparkTask
//
//	@Description: 导入流水到图数据库, Mongo数据库
//	@receiver s
//	@param ctx
//	@param caseId
//	@param taskId
//	@param s3Path 要导入的s3文件
//	@param hugeGraphName 图名称
//	@return error
func (s *RecordService) CreateSparkTask(ctx context.Context, cfg CreateSparkTaskOptions) (err error) {
	defer func() {
		if err == nil {
			err = errors.GetError(recover())
		}
	}()
	tenantId := appctx.GetTenantId2(ctx)
	// 1. 构建 S3 认证与 Spark 运行配置
	sparkConf := map[string]string{
		"spark.hadoop.fs.s3a.impl":                              "org.apache.hadoop.fs.s3a.S3AFileSystem", // 核心 S3A 文件系统驱动
		"spark.hadoop.fs.s3a.endpoint":                          cfg.S3.Endpoint,                          // "http://192.168.120.224:9000",            // 如果是自建的 MinIO、Ceph 或对象存储，填入具体的 "IP:端口" 或 "域名"
		"spark.hadoop.fs.s3a.access.key":                        cfg.S3.User,                              // "minioadmin",                             // 设置登录凭证 (Access Key 和 Secret Key)
		"spark.hadoop.fs.s3a.secret.key":                        cfg.S3.Password,                          // "minioadmin",                             //
		"spark.hadoop.fs.s3a.path.style.access":                 "true",                                   // 自建 S3 服务（如 MinIO）的两个非常重要的避坑配置： 必须开启 path style access，否则 Spark 会尝试访问 bucket.192.168.1.100 导致 DNS 解析失败
		"spark.hadoop.fs.s3a.connection.ssl.enabled":            "false",                                  // 如果你的内部 S3 服务使用的是 HTTP 而不是 HTTPS，必须将其设为 false
		"spark.sql.streaming.forceDeleteTempCheckpointLocation": "true",                                   //
		"spark.driver.userClassPathFirst":                       "false",                                  //
		"spark.executor.userClassPathFirst":                     "false",                                  //
		"spark.executor.heartbeatInterval":                      "60s",                                    // 1. 将 Executor 向 Driver 汇报心跳的间隔从 10s 延长到 60s
		"spark.network.timeout":                                 "800s",                                   // 2. 将全局网络超时时间从 120s 延长到 800s (必须大于 heartbeatInterval)
		"spark.task.maxFailures":                                "4",                                      // 3. 增加重试次数，防止偶尔的网络抖动导致写入失败
		"spark.master":                                          "local[*]",
		//"spark.driver.extraJavaOptions":                         "-Dhttp.proxyHost= -Dhttp.proxyPort= -Dhttps.proxyHost= -Dhttps.proxyPort=", //
		//"spark.executor.extraJavaOptions":                       "-Dhttp.proxyHost= -Dhttp.proxyPort= -Dhttps.proxyHost= -Dhttps.proxyPort=", //
	}

	/*
		--app-id records.import2doris
		--debug true
		--s3-path s3a://import-data/record/record.parquet
		--task-id 2ndoebXy02tqJTdy9YisXp4I
		--tenant-id test
		--case-id 1001
		--host 192.168.120.224
		--port 8030
		--user root
		--db master
		--table master_record
	*/

	// 2. 构建 Spark 任务参数
	request := livyserver.CreateBatchRequest{
		File:      cfg.JavaJarFile,
		ClassName: cfg.JavaApplicationClassName, //"org.example.SparkJavaApp",
		Args: []string{
			"--app-id", "records.import2doris", // spark appid
			"--debug", "true",
			"--tenant-id", tenantId, // tenant id
			"--task-id", cfg.Args.TaskId, // task id
			"--case-id", cfg.Args.CaseId, // case id
			"--s3-data-file", cfg.Args.S3DataFile, // s3 data file
			"--host", cfg.Args.Doris.Host, // "192.168.120.200",
			"--port", cfg.Args.Doris.Port, // "18080",
			"--user", cfg.Args.Doris.User, //"admin",
			"--pwd", cfg.Args.Doris.Password, // "admin",
		},
		Conf:           sparkConf,
		DriverMemory:   cfg.DriverMemory,
		ExecutorMemory: cfg.ExecutorMemory,
		ExecutorCores:  cfg.ExecutorCores,
	}
	_, err = s.livyServer.CreateBatch(request)
	return err
}
