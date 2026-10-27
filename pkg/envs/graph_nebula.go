package envs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const defaultNebulaPort = 9669

type Nebula struct {
	DbKey    string                 `yaml:"dbKey" json:"dbKey"`
	Name     string                 `yaml:"-"`
	Addrs    []string               `yaml:"addrs"`
	User     string                 `yaml:"user"`
	Password string                 `yaml:"password"`
	Space    string                 `yaml:"space"`
	PoolSize int                    `yaml:"poolSize"`
	pool     *nebula.ConnectionPool `yaml:"-"`
}

func NewNebulaWithMap(meta map[string]any) (*Nebula, error) {
	var val Nebula
	if err := copyObject(meta, &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func AddGraphDBNebula(env *Env, cfg *Nebula) error {
	if len(cfg.Addrs) == 0 {
		return fmt.Errorf("neo4j address is empty")
	}
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}
	addrs := toNebulaHostAddresses(cfg.Addrs)
	pool, err := nebula.NewConnectionPool(addrs, nebula.PoolConfig{
		MaxConnPoolSize: cfg.PoolSize,
	}, nebula.DefaultLogger{})

	if err != nil {
		panic(fmt.Sprintf("nebula 连接池创建失败: %v", err))
	}

	cfg.pool = pool
	dbKey := strings.ToLower(cfg.DbKey)
	cfg.Name = dbKey
	env.AddDB(&dbItem{dbKey: dbKey, dbType: DBTypeMongoDB, nebulaPool: pool, config: cfg})
	return nil
}

func (env *Env) GetNebulaClient(dbKey string) *nebula.ConnectionPool {
	item := env.GetDB(dbKey)
	if item == nil {
		panic(fmt.Sprintf("dbKey %s not ", dbKey))
	}
	if item.GetDBType() != DBTypeNebula {
		panic(fmt.Sprintf("dbKey %s is not Nebula", dbKey))
	}
	return item.GetNebulaClient()
}

// GetSession 从连接池获取一个已认证的 Nebula Session
func (n *Nebula) GetSession(ctx context.Context) (*nebula.Session, error) {
	if n == nil || n.pool == nil {
		return nil, fmt.Errorf("nebula connection pool not initialized for %q", n.Name)
	}
	return n.pool.GetSession(n.User, n.Password)
}

func toNebulaHostAddresses(addrs []string) []nebula.HostAddress {
	result := make([]nebula.HostAddress, 0, len(addrs))
	for _, addr := range addrs {
		host, port := splitHostPort(addr, defaultNebulaPort)
		result = append(result, nebula.HostAddress{Host: host, Port: port})
	}
	return result
}

func splitHostPort(addr string, defaultPort int) (string, int) {
	idx := strings.LastIndex(addr, ":")
	if idx < 0 {
		return addr, defaultPort
	}
	host := addr[:idx]
	portStr := addr[idx+1:]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return addr, defaultPort
	}
	return host, port
}
