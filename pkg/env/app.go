package env

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
)

// App
// @Description:  应用配置
// @Author:       liuxd
// @Date:         2021/10/18 10:57
type App struct {
	AppId      string         `yaml:"id" json:"id"`               // 应用ID
	AppName    string         `yaml:"name" json:"name"`           // 应用名称
	ProdMode   bool           `yaml:"prodMode" json:"prodMode"`   // 是生产模式
	HttpHost   string         `yaml:"httpHost" json:"httpHost"`   // 绑定HTTP IP
	HttpPort   int            `yaml:"httpPort" json:"httpPort"`   // 绑定HTTP 端口
	RootUrl    string         `yaml:"rootUrl" json:"rootUrl"`     // URL根
	CPU        *int           `yaml:"cpu" json:"cpu"`             // CPU数量
	Memory     *string        `yaml:"memory" json:"memory"`       // 内存大小
	Meta       map[string]any `yaml:"meta" json:"meta"`           // 自定义配置
	AuthToken  string         `yaml:"authToken" json:"authToken"` // 开发时Token
	HServer    HServer        `yaml:"hServer" json:"hServer"`     // 脚本服务配置
	IsPubEvent bool           `yaml:"isPubEvent" json:"isPubEvent"`
}

// Log 日志配置
type Log struct {
	SaveDays   int           `yaml:"saveDays" json:"saveDays"`   //日志保存的天数
	SplitHour  int           `yaml:"splitHour" json:"splitHour"` //文件分隔时间，单位小时
	LogFile    string        `yaml:"logFile" json:"logFile"`
	OutputType LogOutputType `yaml:"outputType" json:"outputType"` // 日志输出类型 console、 file、 all
	Level      string        `yaml:"level" json:"level"`           // 日志输出类型 console、 file、 all
	LogLevel   logs.Level
}

type Auth struct {
	StoreDbKey string `json:"storeDbKey"`
	Enable     *bool  `json:"enable"`
}

// LogOutputType 日志输出类型
type LogOutputType string

const (
	LogOutputType_Console LogOutputType = "console"
	LogOutputType_File    LogOutputType = "file"
	LogOutputType_All     LogOutputType = "all"
)

func NewLog() *Log {
	return &Log{
		SaveDays:   10,
		SplitHour:  12,
		OutputType: LogOutputType_Console,
		Level:      "info",
		LogLevel:   logs.InfoLevel,
	}
}

func NewDapr() *Dapr {
	return &Dapr{
		Enable: false,
	}
}

func NewApp() *App {
	return &App{
		AppId:    "",
		AppName:  "",
		HttpHost: "0.0.0.0",
		HttpPort: 1984,
	}
}

func initApp(env *Env) {
	if env == nil {
		return
	}
	if env.App == nil {
		env.App = NewApp()
		return
	}
	if env.App.AppId == "" {
		panic("app id is required")
	}

	if len(env.App.HttpHost) == 0 {
		env.App.HttpHost = "0.0.0.0"
	}

	ddd.Init(env.App.AppId)

}

func initLog(e *Env) {
	if e == nil {
		return
	}
	if e.Log == nil {
		e.Log = NewLog()
	}

	// init log
	level := logs.ErrorLevel
	if e.Log.Level != "" {
		l, err := logs.ParseLevel(e.Log.Level)
		if err != nil {
			panic(err)
		}
		level = l
	}
	if e.Log.SaveDays <= 0 {
		e.Log.SaveDays = 30
	}
	if e.Log.SplitHour <= 0 {
		e.Log.SplitHour = 24
	}
	if e.Log.LogFile == "" {
		logFile := fmt.Sprintf("./logs/%s.log", GetExeName())
		e.Log.LogFile = AbsFileName(logFile)
	}
	e.Log.LogLevel = level

}

func NewAuth() *Auth {
	return &Auth{}
}
