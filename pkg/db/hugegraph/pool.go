package hugegraph

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/hgtransport"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

var (
	ErrPoolExhausted  = errors.New("hugegraph connection pool exhausted")
	ErrPoolClosed     = errors.New("hugegraph connection pool is closed")
	ErrTenantNotFound = errors.New("tenant graph pool not found")
)

// PoolConfig 单个租户连接池的配置
type PoolConfig struct {
	InitialCap  int           // 初始连接数
	MaxCap      int           // 最大连接数
	MaxIdle     int           // 最大空闲连接数
	IdleTimeout time.Duration // 空闲连接超时时间
	WaitTimeout time.Duration // 获取连接的最大等待时间
}

// TenantConfig 租户的图实例配置
type HugeConfig struct {
	/*	TenantID string
		Graph    string*/
	Host     string
	Port     int
	Username string
	Password string
}

// hgConn 包装实际的连接和元数据
type hgConn struct {
	client     *hugegraph.CommonClient
	lastActive time.Time
}

// ClientPool 单个租户的图连接池
type ClientPool struct {
	config     PoolConfig
	hugeConfig HugeConfig
	mu         sync.Mutex
	conns      chan *hgConn
	activeCnt  int
	isClosed   bool
}

func GetPoolConfigDefault() PoolConfig {
	return PoolConfig{
		InitialCap:  1,
		MaxCap:      10,
		MaxIdle:     10,
		IdleTimeout: 30 * time.Second,
		WaitTimeout: 30 * time.Second,
	}
}

// NewClientPool 创建单租户连接池
func NewClientPool(hugeCfg HugeConfig, pCfg PoolConfig) (*ClientPool, error) {
	if pCfg.InitialCap > pCfg.MaxCap || pCfg.MaxIdle > pCfg.MaxCap {
		return nil, fmt.Errorf("invalid pool configuration bounds")
	}

	p := &ClientPool{
		config:     pCfg,
		hugeConfig: hugeCfg,
		conns:      make(chan *hgConn, pCfg.MaxCap),
		activeCnt:  0,
		isClosed:   false,
	}

	// 预温连接池
	for i := 0; i < pCfg.InitialCap; i++ {
		client, err := p.newClient()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("failed to init client: %w", err)
		}
		p.conns <- &hgConn{client: client, lastActive: time.Now()}
		p.activeCnt++
	}

	// 启动定时空闲连接清理
	go p.idleCleaner()

	return p, nil
}

// newClient 创建底层的 HugeGraph 客户端
func (p *ClientPool) newClient() (*hugegraph.CommonClient, error) {
	cfg := hugegraph.Config{
		Host:     p.hugeConfig.Host,
		Port:     p.hugeConfig.Port,
		Username: p.hugeConfig.Username,
		Password: p.hugeConfig.Password,
		//Timeout:  5 * time.Second,
		Logger: &hgtransport.ColorLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	}
	return hugegraph.NewCommonClient(cfg)
}

// Get 获取一个连接
func (p *ClientPool) Get(ctx context.Context) (*hugegraph.CommonClient, error) {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return nil, ErrPoolClosed
	}

	// 1. 尝试从已有空闲通道获取
	select {
	case conn := <-p.conns:
		p.mu.Unlock()
		// 检查连接是否因空闲过久失效，或进行 Ping 健康检查
		if time.Since(conn.lastActive) > p.config.IdleTimeout {
			clientClose(conn.client) // 假定 client 有 Close 方法
			// 失效则重新创建一个
			return p.recreateClient()
		}
		return conn.client, nil
	default:
		// 2. 通道无空闲，检查是否可以创建新连接
		if p.activeCnt < p.config.MaxCap {
			p.activeCnt++
			p.mu.Unlock()
			client, err := p.newClient()
			if err != nil {
				p.mu.Lock()
				p.activeCnt--
				p.mu.Unlock()
				return nil, err
			}
			return client, nil
		}
		p.mu.Unlock()
	}

	// 3. 达到最大连接限制，阻塞等待或直到上下文超时
	waitCtx, cancel := context.WithTimeout(ctx, p.config.WaitTimeout)
	defer cancel()

	select {
	case conn := <-p.conns:
		return conn.client, nil
	case <-waitCtx.Done():
		return nil, ErrPoolExhausted
	}
}

// Put 释放连接归还给池
func (p *ClientPool) Put(client *hugegraph.CommonClient) {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		clientClose(client)
		return
	}

	// 如果当前空闲队列未满，则放回
	select {
	case p.conns <- &hgConn{client: client, lastActive: time.Now()}:
		p.mu.Unlock()
	default:
		// 空闲队列慢了（超过 MaxIdle 但在 MaxCap 内），直接销毁该连接
		p.activeCnt--
		p.mu.Unlock()
		clientClose(client)
	}
}

func clientClose(client *hugegraph.CommonClient) {
	if client != nil {
		//client.Close()
	}
}

func (p *ClientPool) recreateClient() (*hugegraph.CommonClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	client, err := p.newClient()
	if err != nil {
		p.activeCnt--
		return nil, err
	}
	return client, nil
}

// idleCleaner 定时清理超时空闲连接
func (p *ClientPool) idleCleaner() {
	ticker := time.NewTicker(p.config.IdleTimeout / 2)
	for range ticker.C {
		p.mu.Lock()
		if p.isClosed {
			p.mu.Unlock()
			ticker.Stop()
			return
		}

		connsLen := len(p.conns)
		// 保证至少留下 InitialCap 的连接，其余超时的释放
		for i := 0; i < connsLen; i++ {
			select {
			case conn := <-p.conns:
				if time.Since(conn.lastActive) > p.config.IdleTimeout && p.activeCnt > p.config.InitialCap {
					p.activeCnt--
					clientClose(conn.client)
				} else {
					// 未超时或不能再释放，放回队列
					p.conns <- conn
				}
			default:
				break
			}
		}
		p.mu.Unlock()
	}
}

// Close 关闭连接池
func (p *ClientPool) Close() {
	p.mu.Lock()
	if p.isClosed {
		p.mu.Unlock()
		return
	}
	p.isClosed = true
	close(p.conns)

	for conn := range p.conns {
		if conn.client != nil {
			clientClose(conn.client)
		}
	}
	p.activeCnt = 0
	p.mu.Unlock()
}
