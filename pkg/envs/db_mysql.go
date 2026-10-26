package envs

import (
	"fmt"
	"strings"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySql struct {
	DbKey        string          `yaml:"dbKey" json:"dbKey"`
	User         string          `yaml:"user" json:"user"`                 // MySQL 用户名
	Password     string          `yaml:"password" json:"password"`         // MySQL 密码
	Host         string          `yaml:"host" json:"host"`                 // MySQL 主机地址
	Port         string          `yaml:"port" json:"port"`                 // MySQL 端口
	DbName       string          `yaml:"dbname" json:"dbName"`             // 数据库名称
	Charset      string          `yaml:"charset" json:"charset"`           // 字符集
	ParseTime    *bool           `yaml:"parseTime" json:"parseTime"`       // 是否解析时间
	Loc          string          `yaml:"loc" json:"loc"`                   // 时区
	LogLevel     logger.LogLevel `yaml:"loglevel" json:"logLevel"`         // 日志级别
	EventPublish bool            `yaml:"eventPublish" json:"eventPublish"` // 是否发送领域事件
}

func NewMySQL() *MySql {
	return &MySql{}
}

func NewMySQLWithMap(meta map[string]any) (*MySql, error) {
	var val MySql
	if err := copyObject(meta, &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func AddDBMySql(env *Env, cfg *MySql) error {
	if cfg.Host == "<no value>" && cfg.Port == "<no value>" {
		return errors.New("")
	}
	logLevel := getGormLogLevel(env.Log.LogLevel)
	dsn := cfg.DSN()
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  false, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: true,  // 根据当前 MySQL 版本自动配置
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		},
	)
	if err != nil {
		panic(fmt.Sprintf("%s ; 连接mysql失败, error:%s", dsn, err.Error()))
	}

	item := &dbItem{
		dbKey:  cfg.DbKey,
		dbType: DBTypeMySQL,
		gormDb: db,
		config: cfg,
	}
	env.AddDB(item)
	return nil
}

// DSN 构造 MySQL 连接字符串
func (cfg *MySql) DSN() string {
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

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v&loc=%s&interpolateParams=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, charset, parseTime, loc)
	return dsn
}

func (cfg *MySql) params(val string) string {
	val = strings.ReplaceAll(val, "&", "%26")
	val = strings.ReplaceAll(val, "/", "%2F")
	val = strings.ReplaceAll(val, "=", "%3D")
	return val
}

func getGormLogLevel(level logs.Level) logger.LogLevel {
	switch level {
	case logs.TraceLevel:
	case logs.DebugLevel, logs.InfoLevel:
		return logger.Info
	case logs.WarnLevel:
		return logger.Warn
	default:
		return logger.Info
	}
	return logger.Info
}
