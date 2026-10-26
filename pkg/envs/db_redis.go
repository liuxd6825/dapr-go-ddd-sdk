package envs

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	DbKey           string  `yaml:"dbKey" json:"dbKey"`
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
	Client          *redis.Client
}

func NewRedis() *Redis {
	return &Redis{}
}

func NewRedisWithMap(meta map[string]any) (*Redis, error) {
	var val Redis
	if err := copyObject(meta, &val); err != nil {
		return nil, err
	}
	return &val, nil
}

func AddDBRedis(env *Env, cfg *Redis) error {
	second5 := 5 * time.Second
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Password, // no password set
		DB:       cfg.Database, // use default DB
		// 连接池最大socket连接数，默认为4倍CPU数
		PoolSize: parseInt("redis.poolSize", cfg.PoolSize, 10),
		//读超时，默认3秒， -1表示取消读超时
		ReadTimeout: parseDuration("redis.readTimeout", cfg.ReadTimeout, second5),
		//写超时，默认等于读超时
		WriteTimeout: parseDuration("redis.writeTimeout", cfg.WriteTimeout, second5),
		//连接建立超时时间，默认5秒。
		DialTimeout: parseDuration("redis.dialTimeout", cfg.DialTimeout, second5),
		//当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒
		PoolTimeout: parseDuration("redis.poolTimeout", cfg.PoolTimeout, second5),
		// 命令执行失败时，最多重试多少次，默认为0即不重试
		MaxRetries: parseInt("redis.maxRetries", cfg.MaxRetries, 3),
		// 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MinRetryBackoff: parseDuration("redis.minRetryBackoff", cfg.MinRetryBackoff, 8*time.Millisecond),
		// 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
		MaxRetryBackoff: parseDuration("redis.maxRetryBackoff", cfg.MaxRetryBackoff, 512*time.Millisecond),
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logs.Errorf(context.Background(), nil, "连接neo4j失败, host:%s, error:%s  ", cfg.Host, err.Error())
		os.Exit(0)
	}
	cfg.DbKey = strings.ToLower(cfg.DbKey)
	cfg.Client = rdb
	item := &dbItem{
		dbKey:  cfg.DbKey,
		dbType: DBTypeRedis,
		redis:  rdb,
		config: cfg,
	}
	env.AddDB(item)
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
