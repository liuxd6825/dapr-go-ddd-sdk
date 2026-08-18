package env

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

const defaultNebulaPort = 9669

type Nebula struct {
	Name     string                    `yaml:"-"`
	Addrs    []string                  `yaml:"addrs"`
	User     string                    `yaml:"user"`
	Password string                    `yaml:"password"`
	Space    string                    `yaml:"space"`
	PoolSize int                       `yaml:"poolSize"`
	pool     *nebula.ConnectionPool    `yaml:"-"`
}

func InitDBNebula(env *Env) {
	if env.Nebula == nil {
		env.Nebula = map[string]*Nebula{}
		return
	}
	for dbKey, cfg := range env.Nebula {
		if len(cfg.Addrs) == 0 {
			continue
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
		key := strings.ToLower(dbKey)
		cfg.Name = key
	}
}

func (env *Env) GetNebulaByKey(dbKey string) (*Nebula, bool) {
	if env.Nebula == nil {
		return nil, false
	}
	n := env.Nebula[dbKey]
	return n, n != nil
}

func (env *Env) AddNebula(cfg *Nebula) {
	if cfg == nil {
		return
	}
	if env.Nebula == nil {
		env.Nebula = map[string]*Nebula{}
	}
	env.Nebula[cfg.Name] = cfg
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