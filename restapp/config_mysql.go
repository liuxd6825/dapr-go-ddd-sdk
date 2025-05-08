package restapp

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
)

type MySqlConfig struct {
	Host     string          `yaml:"host"`
	Port     string          `yaml:"port"`
	Database string          `yaml:"dbname"`
	UserName string          `yaml:"user"`
	Password string          `yaml:"pwd"`
	LogLevel logger.LogLevel `yaml:"loglevel"`
}

var _mysqlList map[string]*gorm.DB
var _mysqlDefault *gorm.DB

func initMySql(configs map[string]*MySqlConfig) {
	if err := assert.NotNil(configs, assert.NewOptions("cfg is nil")); err != nil {
		panic(err)
	}

	for key, cfg := range configs {
		if cfg.Host == "<no value>" && cfg.Port == "<no value>" {
			continue
		}
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local", cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,   // DSN data source name
			DefaultStringSize:         256,   // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}), &gorm.Config{Logger: logger.Default.LogMode(cfg.LogLevel)})

		if err != nil {
			logs.Errorf(context.Background(), "", nil, "连接mysql失败, error:%s", err.Error())
			os.Exit(0)
		}
		_mysqlList[strings.ToLower(key)] = db
		_mysqlDefault = db
	}
}

func GetMySql() *gorm.DB {
	return _mysqlDefault
}

func GetMySqlByKey(dbKey string) (*gorm.DB, bool) {
	if len(dbKey) == 0 {
		return _mysqlDefault, _mysqlDefault != nil
	}
	d, ok := _mysqlList[strings.ToLower(dbKey)]
	return d, ok
}

func CloseAllMySql(ctx context.Context) error {
	c := func(d *gorm.DB) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		return err
	}
	for _, d := range _mysqlList {
		_ = c(d)
	}
	return nil
}
