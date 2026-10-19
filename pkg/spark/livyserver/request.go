package livyserver

// CreateBatchRequest 定义提交任务的请求体
type CreateBatchRequest struct {
	File           string            `json:"file"`           // S3 JAR 路径 (s3a://...)
	ClassName      string            `json:"className"`      // 主类入口
	Args           []string          `json:"args,omitempty"` // 传给 main 方法的参数
	Jars           []string          `json:"jars,omitempty"`
	Conf           map[string]string `json:"conf,omitempty"` // Spark / Hadoop S3 配置
	DriverMemory   string            `json:"driverMemory,omitempty"`
	ExecutorMemory string            `json:"executorMemory,omitempty"`
	ExecutorCores  int               `json:"executorCores,omitempty"`
}
