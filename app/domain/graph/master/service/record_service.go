package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
	envs "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
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

func NewRecordServiceWithEnv(env *envs.Env, graphCfg graphconfigs.Config) *RecordService {
	livyServer := livyserver.NewLivyServer(graphCfg.LivyServer.IP, graphCfg.LivyServer.Port)
	s3cfg := env.GetMinioConfig(graphCfg.OssKey)
	hgCfg := env.GetHugeConfig(graphCfg.GraphDbKey)
	mongoCfg := env.GetMongoConfig(graphCfg.MongoDbKey)

	cfg := RecordServiceConfig{
		S3Endpoint: s3cfg.Endpoint,
		S3User:     s3cfg.AccessKey,
		S3Pwd:      s3cfg.SecretKey,

		HGHost: hgCfg.Host,
		HGPort: hgCfg.Port,
		HGUser: hgCfg.Username,
		HGPwd:  hgCfg.Password,

		MongoHost:       mongoCfg.Host,
		MongoReplicaset: mongoCfg.ReplicaSet,
		MongoUser:       mongoCfg.User,
		MongoPwd:        mongoCfg.Password,
		MongoDB:         mongoCfg.DbName,
	}

	return NewRecordService(livyServer, cfg)
}

type RecordServiceConfig struct {
	// S3 要导入的文件
	S3Endpoint string // "http://192.168.120.224:9000"
	S3User     string //
	S3Pwd      string //

	// 图数据库
	HGHost  string // "192.168.120.200"
	HGPort  int    // "18080"
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

type CreateImportTaskOptions struct {
	DriverMemory   string // "8g"
	ExecutorMemory string // "8g"
	ExecutorCores  int    //  6
}

func NewCreateImportTaskOptions(opts ...CreateImportTaskOptions) *CreateImportTaskOptions {
	o := &CreateImportTaskOptions{}
	for _, i := range opts {
		if i.DriverMemory != "" {
			o.DriverMemory = i.DriverMemory
		}
		if i.ExecutorMemory != "" {
			o.ExecutorMemory = i.ExecutorMemory
		}
		if i.ExecutorCores != 0 {
			o.ExecutorCores = i.ExecutorCores
		}
	}
	return o
}

// CreateImportTask
//
//	@Description: 导入流水到图数据库, Mongo数据库
//	@receiver s
//	@param ctx
//	@param caseId
//	@param taskId
//	@param s3Path 要导入的s3文件
//	@param hugeGraphName 图名称
//	@return error
func (s *RecordService) CreateImportTask(ctx context.Context, caseId, taskId, s3DataFile, s3JarFile, hugeGraphName string, opts ...CreateImportTaskOptions) (sparkTaskId string, err error) {
	opt := NewCreateImportTaskOptions(opts...)
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
	request := livyserver.CreateBatchRequest{
		File:      s3JarFile,
		ClassName: "org.example.SparkJavaApp",
		Args: []string{
			"--tenant-id", tenantId, //
			"--task-id", taskId, //
			"--case-id", caseId, //
			"--s3-path", s3DataFile, //
			"--hg-host", s.cfg.HGHost, // "192.168.120.200",
			"--hg-port", fmt.Sprintf("%d", s.cfg.HGPort), // "18080",
			"--hg-graph", hugeGraphName, // "hugegraph",
			"--hg-user", s.cfg.HGUser, //"admin",
			"--hg-pwd", s.cfg.HGPwd, // "admin",
			"--mongo-host", s.cfg.MongoHost, // "192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020",
			"--mongo-replicaset", s.cfg.MongoReplicaset, // "mongors",
			"--mongo-db", s.cfg.MongoDB, //"master",
			"--mongo-user", s.cfg.MongoUser, // "super_admin",
			"--mongo-pwd", s.cfg.MongoPwd, //"123456",
		},
		Conf:           sparkConf,
		DriverMemory:   opt.DriverMemory,
		ExecutorMemory: opt.ExecutorMemory,
		ExecutorCores:  opt.ExecutorCores,
	}

	sparkTask, err := s.livyServer.CreateBatch(request)
	if err != nil {
		return "", err
	}

	if sparkTask == nil {
		return "", errors.New("spark task is nil")
	}
	return strconv.Itoa(sparkTask.ID), nil
}
