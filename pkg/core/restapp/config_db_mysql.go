package restapp

import (
	"fmt"
	"strings"

	"gorm.io/gorm/logger"
)

// MySqlConfig 结构体用于存储 MySQL 连接信息
type MySqlConfig struct {
	User         string          `yaml:"user" json:"username"`             // MySQL 用户名
	Password     string          `yaml:"pwd" json:"password"`              // MySQL 密码
	Host         string          `yaml:"host" json:"host"`                 // MySQL 主机地址
	Port         string          `yaml:"port" json:"port"`                 // MySQL 端口
	DbName       string          `yaml:"dbname" json:"dbName"`             // 数据库名称
	Charset      string          `yaml:"charset" json:"charset"`           // 字符集
	ParseTime    *bool           `yaml:"parseTime" json:"parseTime"`       // 是否解析时间
	Loc          string          `yaml:"loc" json:"loc"`                   // 时区
	LogLevel     logger.LogLevel `yaml:"loglevel" json:"logLevel"`         // 日志级别
	EventPublish bool            `yaml:"eventPublish" json:"eventPublish"` // 是否发送领域事件
}

// DSN 构造 MySQL 连接字符串
func (cfg *MySqlConfig) DSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

	charset := "utf8"
	if cfg.Charset != "" {
		charset = cfg.params(cfg.Charset)
	}

	parseTime := true
	if cfg.ParseTime != nil {
		parseTime = *cfg.ParseTime
	}

	loc := "Local"
	if cfg.Loc != "" {
		loc = cfg.params(cfg.Loc)
	}

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v&loc=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, charset, parseTime, loc)
	return dsn
}

func (cfg *MySqlConfig) params(val string) string {
	val = strings.ReplaceAll(val, "&", "%26")
	val = strings.ReplaceAll(val, "/", "%2F")
	val = strings.ReplaceAll(val, "=", "%3D")
	return val
}
