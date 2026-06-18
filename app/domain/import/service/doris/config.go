package doris

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"
)

// Config 存储 Doris 连接配置
type Config struct {
	BEHost   string        // BE 节点 IP 或域名
	BEPort   int           // BE HTTP 端口 (默认 8040)
	FEHost   string        // FE 节点 IP 或域名,留空时回退到 BEHost
	FEPort   int           // FE MySQL 端口 (默认 9030)
	Database string        // 数据库名
	Table    string        // 表名
	Username string        // 账号
	Password string        // 密码
	Timeout  time.Duration // 请求超时时间
}

func NewConfigWithMap(val map[string]any) (*Config, error) {
	if val == nil {
		return nil, errors.New("config is nil")
	}
	config := &Config{}
	if host, err := maputils.GetString(val, "host", ""); err != nil {
		return nil, err
	} else {
		config.BEHost = host
	}
	if port, err := maputils.GetInt(val, "port", 0); err != nil {
		return nil, err
	} else {
		config.BEPort = port
	}
	if database, err := maputils.GetString(val, "database", ""); err != nil {
		return nil, err
	} else {
		config.Database = database
	}
	if username, err := maputils.GetString(val, "username", ""); err != nil {
		return nil, err
	} else {
		config.Username = username
	}
	if password, err := maputils.GetString(val, "password", ""); err != nil {
		return nil, err
	} else {
		config.Password = password
	}
	if timeout, err := maputils.GetInt(val, "timeout", 3*1000); err != nil {
		return nil, err
	} else {
		config.Timeout = time.Duration(timeout) * time.Millisecond
	}
	return config, nil
}
