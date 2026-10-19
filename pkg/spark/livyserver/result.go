package livyserver

// CreateTaskResult 定义 Livy 返回的响应体
type CreateTaskResult struct {
	ID      int                    `json:"id"`
	State   string                 `json:"state"` // starting, running, success, dead, killed
	AppID   string                 `json:"appId"`
	AppInfo map[string]interface{} `json:"appInfo"`
	Log     []string               `json:"log"`
}

type GetBatchResult struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Owner     string `json:"owner"`
	ProxyUser string `json:"proxyUser"`
	State     string `json:"state"`
	AppId     string `json:"appId"`
	AppInfo   struct {
		DriverLogUrl    string `json:"driverLogUrl"`
		SparkUiUrl      string `json:"sparkUiUrl"`
		ExecutorLogUrls string `json:"executorLogUrls"`
	} `json:"appInfo"`
	Log []string `json:"log"`
}
type GetBatchLogsResult struct {
}

type GetBatchStateResult struct {
	Id    string `json:"id"`
	State string `json:"state"`
}

type GetSessionResult struct {
	From     int       `json:"from"`
	Total    int       `json:"total"`
	Sessions []Session `json:"session"`
}
