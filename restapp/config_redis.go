package restapp

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"time"
)

type RedisConfig struct {
	DbKey           string
	Host            string  `yaml:"host"`
	Database        int     `yaml:"db"`
	Password        string  `yaml:"pwd"`
	PoolSize        *int    `yaml:"poolSize"`
	ReadTimeout     *string `yaml:"readTimeout"`
	WriteTimeout    *string `yaml:"writeTimeout"`
	DialTimeout     *string `yaml:"dialTimeout"`
	PoolTimeout     *string `yaml:"poolTimeout"`
	MaxRetries      *int    `yaml:"maxRetries"`
	MinRetryBackoff *string `yaml:"minRetryBackoff"`
	MaxRetryBackoff *string `yaml:"maxRetryBackoff"`
}

var _redisDbs map[string]*redis.Client
var _initRedis bool = false
var _redisDefault *redis.Client

func init() {
	_redisDbs = make(map[string]*redis.Client)
}

func initRedis(configs map[string]*RedisConfig) error {
	if _initRedis {
		return nil
	}
	_initRedis = true
	if err := assert.NotNil(configs, assert.NewOptions("configs is nil")); err != nil {
		panic(err)
	}
	second5 := 5 * time.Second
	for k, c := range configs {
		rdb := redis.NewClient(&redis.Options{
			Addr:     c.Host,
			Password: c.Password, // no password set
			DB:       c.Database, // use default DB
			// 连接池最大socket连接数，默认为4倍CPU数
			PoolSize: parseInt("redis.poolSize", c.PoolSize, 10),
			//读超时，默认3秒， -1表示取消读超时
			ReadTimeout: parseDuration("redis.readTimeout", c.ReadTimeout, second5),
			//写超时，默认等于读超时
			WriteTimeout: parseDuration("redis.writeTimeout", c.WriteTimeout, second5),
			//连接建立超时时间，默认5秒。
			DialTimeout: parseDuration("redis.dialTimeout", c.DialTimeout, second5),
			//当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
			PoolTimeout: parseDuration("redis.poolTimeout", c.PoolTimeout, second5),
			// 命令执行失败时，最多重试多少次，默认为0即不重试
			MaxRetries: parseInt("redis.maxRetries", c.MaxRetries, 3),
			// 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
			MinRetryBackoff: parseDuration("redis.minRetryBackoff", c.MinRetryBackoff, 8*time.Millisecond),
			// 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
			MaxRetryBackoff: parseDuration("redis.maxRetryBackoff", c.MaxRetryBackoff, 512*time.Millisecond),
		})

		if _, err := rdb.Ping(context.Background()).Result(); err != nil {
			logs.Errorf(context.Background(), "", nil, "连接neo4j失败, host:%s, error:%s  ", c.Host, err.Error())
			os.Exit(0)
		}
		dbKey := strings.ToLower(k)
		_redisDbs[dbKey] = rdb
		_redisDefault = rdb
	}
	if len(_redisDbs) > 1 {
		_redisDefault = nil
	}
	return nil
}

func parseInt(name string, val *int, defaultVal int) int {
	if val == nil {
		return defaultVal
	}
	return *val
}

func parseDuration(name string, val *string, defaultDur time.Duration) time.Duration {
	if val == nil || *val == "" {
		return defaultDur
	}
	dur, err := time.ParseDuration(*val)
	if err != nil {
		panic(name + err.Error())
	}
	return dur
}

func GetRedis() *redis.Client {
	if _redisDefault != nil {
		return _redisDefault
	}
	for _, item := range _redisDbs {
		return item
	}
	return nil
}

func GetRedisByKey(dbKey string) (*redis.Client, bool) {
	d, ok := _redisDbs[strings.ToLower(dbKey)]
	return d, ok
}

func SetRedis(dbKey string, client *redis.Client) {
	key := strings.ToLower(dbKey)
	_redisDbs[key] = client
	if _redisDefault == nil {
		_redisDefault = client
	}
}

func CloseRedis(ctx context.Context) error {
	c := func(d *redis.Client) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		return d.Close()
	}
	for _, d := range _redisDbs {
		_ = c(d)
	}
	return nil
}
