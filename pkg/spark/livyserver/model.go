package livyserver

// SessionKind 定义会话类型 (spark, pyspark, sparkr, or sql)
// SessionKind 定义会话类型的专属类型
type SessionKind string

// 定义支持的 Session Kind 常量
const (
	SessionKindSpark   SessionKind = "spark"   // SessionKindSpark 表示交互式 Scala Spark 会话
	SessionKindPySpark SessionKind = "pyspark" // SessionKindPySpark 表示交互式 Python Spark 会话
	SessionKindSparkR  SessionKind = "sparkr"  // SessionKindSparkR 表示交互式 R Spark 会话
	SessionKindSQL     SessionKind = "sql"     // SessionKindSQL 表示交互式 SQL Spark 会话
)

func (k SessionKind) IsValid() bool {
	switch k {
	case SessionKindSpark, SessionKindPySpark, SessionKindSparkR, SessionKindSQL:
		return true
	}
	return false
}

// Session 代表文档中的会话信息结构体
type Session struct {
	ID             int               `json:"id"`             // The session id
	AppID          string            `json:"appId"`          // The application id of this session
	Owner          string            `json:"owner"`          // Remote user who submitted this session
	ProxyUser      string            `json:"proxyUser"`      // User to impersonate when running
	Kind           SessionKind       `json:"kind"`           // Session kind (spark, pyspark, sparkr, or sql)
	Log            []string          `json:"log"`            // The log lines
	State          string            `json:"state"`          // The session state
	AppInfo        map[string]string `json:"appInfo"`        // The detailed application info (Map of key=val)
	Jars           []string          `json:"jars"`           // jars to be used in this session
	PyFiles        []string          `json:"pyFiles"`        // Python files to be used in this session
	Files          []string          `json:"files"`          // files to be used in this session
	DriverMemory   string            `json:"driverMemory"`   // Amount of memory to use for the driver process
	DriverCores    int               `json:"driverCores"`    // Number of cores to use for the driver process
	ExecutorMemory string            `json:"executorMemory"` // Amount of memory to use per executor process
	ExecutorCores  int               `json:"executorCores"`  // Number of cores to use for each executor
	NumExecutors   int               `json:"numExecutors"`   // Number of executors to launch for this session
	Archives       []string          `json:"archives"`       // Archives to be used in this session
	Queue          string            `json:"queue"`          // The name of the YARN queue to which submitted
	Conf           map[string]string `json:"conf"`           // Spark configuration properties (Map of key=val)
}
