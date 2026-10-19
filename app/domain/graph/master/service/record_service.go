package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"
)

// RecordService
// @Description: 银行流水导入到图数据库服务
type RecordService struct {
	livyServer *livyserver.LivyServer
	cfg        RecordServiceConfig
}

func NewRecordService(livyServer *livyserver.LivyServer, cfg RecordServiceConfig) *RecordService {
	return &RecordService{
		livyServer: livyServer,
		cfg:        cfg,
	}
}

type RecordServiceConfig struct {
	// S3 要导入的文件
	S3Endpoint string // "http://192.168.120.224:9000"
	S3User     string //
	S3Pwd      string //

	// 图数据库
	HGHost  string // "192.168.120.200"
	HGPort  string // "18080"
	HGGraph string // "hugegraph"
	HGUser  string
	HGPwd   string

	// 任务状态更新数据库
	MongoHost       string // "192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020"
	MongoReplicaset string // "mongors"
	MongoDB         string // "master"
	MongoUser       string
	MongoPwd        string

	// Spark任务执行
	// LivyURL string // Spark Livy URL "http://192.168.120.200:8998/batches"
}

// Import
//
//	@Description: 导入流水到图数据库, Mongo数据库
//	@receiver s
//	@param ctx
//	@param caseId
//	@param taskId
//	@param s3Path 要导入的s3文件
//	@param hugeGraphName 图名称
//	@return error
func (s *RecordService) Import(ctx context.Context, caseId, taskId, s3DataFile, s3JarFile, hugeGraphName string) error {
	tenantId := appctx.GetTenantId2(ctx)
	// 1. 构建 S3 认证与 Spark 运行配置
	sparkConf := map[string]string{
		"spark.hadoop.fs.s3a.impl":                              "org.apache.hadoop.fs.s3a.S3AFileSystem", // 核心 S3A 文件系统驱动
		"spark.hadoop.fs.s3a.endpoint":                          s.cfg.S3Endpoint,                         // "http://192.168.120.224:9000",            // 如果是自建的 MinIO、Ceph 或对象存储，填入具体的 "IP:端口" 或 "域名"
		"spark.hadoop.fs.s3a.access.key":                        s.cfg.S3User,                             // "minioadmin",                             // 设置登录凭证 (Access Key 和 Secret Key)
		"spark.hadoop.fs.s3a.secret.key":                        s.cfg.S3Pwd,                              // "minioadmin",                             //
		"spark.hadoop.fs.s3a.path.style.access":                 "true",                                   // 自建 S3 服务（如 MinIO）的两个非常重要的避坑配置： 必须开启 path style access，否则 Spark 会尝试访问 bucket.192.168.1.100 导致 DNS 解析失败
		"spark.hadoop.fs.s3a.connection.ssl.enabled":            "false",                                  // 如果你的内部 S3 服务使用的是 HTTP 而不是 HTTPS，必须将其设为 false
		"spark.sql.streaming.forceDeleteTempCheckpointLocation": "true",                                   //
		//"spark.driver.extraJavaOptions":                         "-Dhttp.proxyHost= -Dhttp.proxyPort= -Dhttps.proxyHost= -Dhttps.proxyPort=", //
		//"spark.executor.extraJavaOptions":                       "-Dhttp.proxyHost= -Dhttp.proxyPort= -Dhttps.proxyHost= -Dhttps.proxyPort=", //
		"spark.driver.userClassPathFirst":   "false", //
		"spark.executor.userClassPathFirst": "false", //
		"spark.executor.heartbeatInterval":  "60s",   // 1. 将 Executor 向 Driver 汇报心跳的间隔从 10s 延长到 60s
		"spark.network.timeout":             "800s",  // 2. 将全局网络超时时间从 120s 延长到 800s (必须大于 heartbeatInterval)
		"spark.task.maxFailures":            "4",     // 3. 增加重试次数，防止偶尔的网络抖动导致写入失败
		"spark.master":                      "local[*]",
	}

	// 2. 构建 Spark 任务参数
	request := livyserver.CreateTaskRequest{
		File:      s3JarFile,
		ClassName: "org.example.SparkJavaApp",
		Args: []string{
			"--tenant-id", tenantId, //
			"--task-id", taskId, //
			"--case-id", caseId, //
			"--s3-path", s3DataFile, //
			"--hg-host", s.cfg.HGHost, // "192.168.120.200",
			"--hg-port", s.cfg.HGPort, // "18080",
			"--hg-graph", hugeGraphName, // "hugegraph",
			"--hg-user", s.cfg.HGUser, //"admin",
			"--hg-pwd", s.cfg.HGPwd, // "admin",
			"--mongo-host", s.cfg.MongoHost, // "192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020",
			"--mongo-replicaset", s.cfg.MongoReplicaset, // "mongors",
			"--mongo-db", s.cfg.MongoDB, //"master",
			"--mongo-user", s.cfg.MongoUser, // "super_admin",
			"--mongo-pwd", s.cfg.MongoPwd, //"123456",
		},
		Conf: sparkConf,
		//DriverMemory:   "8g",
		//ExecutorMemory: "8g",
		//ExecutorCores:  6,
	}
	s.livyServer.CreateTasks(ctx, request)
	return nil
}

// trackStatus 根据批处理 ID 轮询任务状态
func (s *RecordService) trackStatus(batchID int) {
	statusURL := s.cfg.LivyURL + "/" + strconv.FormatInt(int64(batchID), 10)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get(statusURL)
		if err != nil {
			fmt.Printf("获取状态失败: %v\n", err)
			continue
		}

		var statusResp LivyBatchResponse
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &statusResp)
		resp.Body.Close()

		fmt.Printf("BatchID [%d] AppID [%s] 当前状态: %s\n", batchID, statusResp.AppID, statusResp.State)
		for _, logInfo := range statusResp.Log {
			fmt.Println(logInfo)
		}
		// 如果状态达到终态则退出轮询
		if statusResp.State == "success" || statusResp.State == "dead" || statusResp.State == "killed" {
			fmt.Printf("任务结束，最终状态: %s\n", statusResp.State)
			break
		}
	}
}
