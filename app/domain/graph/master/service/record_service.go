package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/appctx"
)

// LivyBatchRequest 定义提交 Batch 任务的请求体
type LivyBatchRequest struct {
	File           string            `json:"file"`           // S3 JAR 路径 (s3a://...)
	ClassName      string            `json:"className"`      // 主类入口
	Args           []string          `json:"args,omitempty"` // 传给 main 方法的参数
	Jars           []string          `json:"jars,omitempty"`
	Conf           map[string]string `json:"conf,omitempty"` // Spark / Hadoop S3 配置
	DriverMemory   string            `json:"driverMemory,omitempty"`
	ExecutorMemory string            `json:"executorMemory,omitempty"`
	ExecutorCores  int               `json:"executorCores,omitempty"`
}

// LivyBatchResponse 定义 Livy 返回的响应体
type LivyBatchResponse struct {
	ID      int                    `json:"id"`
	State   string                 `json:"state"` // starting, running, success, dead, killed
	AppID   string                 `json:"appId"`
	AppInfo map[string]interface{} `json:"appInfo"`
	Log     []string               `json:"log"`
}

// RecordService
// @Description: 银行流水导入到图数据库服务
type RecordService struct {
	cfg RecordServiceConfig
}

func NewRecordService(cfg RecordServiceConfig) *RecordService {
	return &RecordService{
		cfg: cfg,
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
	LivyURL string // Spark Livy URL "http://192.168.120.200:8998/batches"
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
	livyURL := s.cfg.LivyURL // "http://192.168.120.200:8998/batches"

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
	requestPayload := LivyBatchRequest{
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

	// 2. 将结构体序列化为 JSON
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		fmt.Printf("JSON 序列化失败: %v\n", err)
		return err
	}

	// 3. 发送 POST 请求到 Livy
	resp, err := http.Post(livyURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("请求 Livy 失败: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// 4. 读取并解析返回结果
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Printf("Livy 报错，状态码: %d, 原因: %s\n", resp.StatusCode, string(body))
		return err
	}

	var livyResp LivyBatchResponse
	if err := json.Unmarshal(body, &livyResp); err != nil {
		fmt.Printf("解析响应 JSON 失败: %v\n", err)
		return err
	}

	fmt.Printf("成功提交 Spark 任务！分配的 Batch ID: %d, 当前状态: %s\n", livyResp.ID, livyResp.State)

	// 5. 异步轮询任务状态（可选流程）
	go s.trackStatus(livyResp.ID)

	// 阻塞主线程方便查看演示结果
	time.Sleep(60 * time.Second)
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
