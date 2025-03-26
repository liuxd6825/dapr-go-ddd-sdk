package restapp

import (
	"context"
	"fmt"
	assert2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
)

// MySqlConfig 结构体用于存储 MySQL 连接信息
type MySqlConfig struct {
	Username     string          `yaml:"user" json:"username"`             // MySQL 用户名
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName)

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

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v&loc=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DbName, charset, parseTime, loc)
	return dsn
}

func (cfg *MySqlConfig) params(val string) string {
	val = strings.ReplaceAll(val, "&", "%26")
	val = strings.ReplaceAll(val, "/", "%2F")
	val = strings.ReplaceAll(val, "=", "%3D")
	return val
}

func initMySql(configs map[string]*MySqlConfig) error {
	if err := assert2.NotNil(configs, assert2.NewOptions("cfg is nil")); err != nil {
		return err
	}

	for key, cfg := range configs {
		if cfg.Host == "<no value>" && cfg.Port == "<no value>" {
			continue
		}
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
				Logger: logger.Default.LogMode(cfg.LogLevel),
			},
		)

		if err != nil {
			logs.Errorf(context.Background(), "", nil, "%s ; 连接mysql失败, error:%s", dsn, err.Error())
			os.Exit(1)
		}

		addMysql(key, db, cfg)
	}
	return nil
}

func addMysql(dbKey string, mysqlDb *gorm.DB, config any) DBItem {
	item := &dbItem{
		dbKey:  dbKey,
		dbType: DbType_MySQL,
		gormDb: mysqlDb,
		config: config,
	}
	addDB(item)
	return item
}
