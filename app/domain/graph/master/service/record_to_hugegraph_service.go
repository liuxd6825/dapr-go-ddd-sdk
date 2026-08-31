package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LivyBatchRequest 定义提交 Batch 任务的请求体
type LivyBatchRequest struct {
	File           string            `json:"file"`           // S3 JAR 路径 (s3a://...)
	ClassName      string            `json:"className"`      // 主类入口
	Args           []string          `json:"args,omitempty"` // 传给 main 方法的参数
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

// Record2HugeGraphService
// @Description: 银行流水 parquet 导入业务编排
type Record2HugeGraphService struct {
}

func NewRecord2HugeGraphService() *RecordParquetService {
	return &RecordParquetService{}
}

// Import
//
//	@Description: 数据导入
//	@receiver s
//	@param ctx 上下文
//	@param s3Path 导入的s3文件地址
//	@param caseId 项目ID
//	@param opts 参数
//	@return error 错误
func (s *RecordParquetService) Import(ctx context.Context, s3Path, caseId string) error {
	livyURL := "http://localhost:8998/batches"

	// 1. 构建 S3 认证与 Spark 运行配置
	sparkConf := map[string]string{
		// 核心 S3A 文件系统驱动
		"spark.hadoop.fs.s3a.impl": "org.apache.hadoop.fs.s3a.S3AFileSystem",
		// 1. 设置 S3 服务的 IP 地址或 Endpoint
		// 如果是自建的 MinIO、Ceph 或对象存储，填入具体的 "IP:端口" 或 "域名"
		"spark.hadoop.fs.s3a.endpoint": "http://192.168.120.224:9000",
		// 2. 设置登录凭证 (Access Key 和 Secret Key)
		"spark.hadoop.fs.s3a.access.key": "minioadmin",
		"spark.hadoop.fs.s3a.secret.key": "minioadmin",
		// 自建 S3 服务（如 MinIO）的两个非常重要的避坑配置：
		// 必须开启 path style access，否则 Spark 会尝试访问 bucket.192.168.1.100 导致 DNS 解析失败
		"spark.hadoop.fs.s3a.path.style.access": "true",
		// 如果你的内部 S3 服务使用的是 HTTP 而不是 HTTPS，必须将其设为 false
		"spark.hadoop.fs.s3a.connection.ssl.enabled":            "false",
		"spark.sql.streaming.forceDeleteTempCheckpointLocation": "true",
	}

	// 1. 构建 Spark 任务参数
	requestPayload := LivyBatchRequest{
		File:      "s3a://spark/app/record-parquet-to-hugegraph-1.0-SNAPSHOT.jar",
		ClassName: "org.example.SparkJavaApp",
		Args: []string{
			"--input-s3-path", "s3a://spark/app/record-parquet-to-hugegraph-1.0-SNAPSHOT.jar",
			"--hg-graph", "hugegraph",
		},
		Conf:           sparkConf,
		DriverMemory:   "3g",
		ExecutorMemory: "5g",
		ExecutorCores:  2,
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
	go trackStatus(livyResp.ID)

	// 阻塞主线程方便查看演示结果
	time.Sleep(60 * time.Second)
	return nil
}

// trackStatus 根据批处理 ID 轮询任务状态
func trackStatus(batchID int) {
	statusURL := fmt.Sprintf("http://localhost:8998/batches/%d", batchID)
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

		fmt.Printf("任务 [%d] 当前状态: %s\n", batchID, statusResp.State)

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
