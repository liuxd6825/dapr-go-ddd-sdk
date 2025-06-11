package env

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type Mongo struct {
	DbKey                  string  `json:"dbKey"`
	AppName                string  `yaml:"appName" json:"appName"`
	Host                   string  `yaml:"host" json:"host"`
	DbName                 string  `yaml:"dbname" json:"dbName"`
	User                   string  `yaml:"user" json:"user"`
	Pwd                    string  `yaml:"pwd" json:"pwd"`
	ReplicaSet             string  `yaml:"replicaSet" json:"replicaSet"`
	WriteConcern           string  `yaml:"writeConcern" json:"writeConcern"`
	ReadConcern            string  `yaml:"readConcern" json:"readConcern"`
	MaxPoolSize            *uint64 `yaml:"maxPoolSize" json:"maxPoolSize"`
	Direct                 *bool   `yaml:"direct" json:"direct" `
	LocalThreshold         string  `yaml:"localThreshold" json:"localThreshold"`                 // 时间长度
	ConnectTimeout         string  `yaml:"connectTimeout" json:"connectTimeout"`                 // 时间长度
	HeartbeatInterval      string  `yaml:"heartbeatInterval" json:"heartbeatInterval"`           // 时间长度
	OperationTimeout       string  `yaml:"operationTimeout" json:"operationTimeout"`             // 时间长度
	MaxConnIdleTime        string  `yaml:"maxConnIdleTime" json:"maxConnIdleTime"`               // 时间长度
	ServerSelectionTimeout string  `yaml:"serverSelectionTimeout" json:"serverSelectionTimeout"` // 时间长度
	SocketTimeout          string  `yaml:"socketTimeout" json:"socketTimeout"`                   // 时间长度
	EventPublish           bool    `yaml:"eventPublish" json:"eventPublish"`                     // 是否发送领域事件
	AutoSource             string  `yaml:"autoSource" json:"autoSource"`                         // 认证源
	AuthMechanism          string  `yaml:"authMechanism" json:"authMechanism"`                   // 认证方式
}

func NewMongo() *Mongo {
	return &Mongo{}
}

func (m *Mongo) IsEmpty() bool {
	if m.Host == "" && m.DbName == "" && m.Pwd == "" && m.User == "" {
		return true
	}
	return false
}

func InitDBMongo(env *Env) {
	if env.Mongo == nil {
		env.Mongo = map[string]*Mongo{}
		return
	}

	for k, c := range env.Mongo {
		if c.IsEmpty() {
			continue
		}
		config := NewStoreMongoConfig(c)
		db, err := mongodb.NewMongoDB(config, func(opts *options.ClientOptions) error {
			logs.Infofmt(context.Background(), "", "config mongo  hosts=%v; user=%s; replicasSet=%s; maxPoolSize=%s; connectTimeout=%v; "+
				"socketTimeout=%v; serverSelectionTimeout=%v; maxConnIdleTime=%v; operationTimeout=%v; socketTimeout=%v ",
				opts.Hosts, opts.Auth.Username, getPStr(opts.ReplicaSet), getPInt(opts.MaxPoolSize), config.ConnectTimeout,
				config.SocketTimeout, config.ServerSelectionTimeout, config.MaxConnIdleTime, config.OperationTimeout, config.SocketTimeout)
			return nil
		})
		if err != nil {
			panic(fmt.Sprintf("mongo连接%s失败, error:%s", config.Host, err.Error()))
		}
		dbKey := strings.ToLower(k)
		c.DbKey = dbKey
		env.AddDB(&dbItem{
			dbKey:  dbKey,
			dbType: DBType_MongoDB,
			mongo:  db,
			config: config,
		})
	}
}

func NewStoreMongoConfig(c *Mongo) *mongodb.Config {
	operationTimeout := getTimeout(c.OperationTimeout, "30s")
	connectTimeout := getTimeout(c.ConnectTimeout, "5s")
	heartbeatInterval := getTimeout(c.HeartbeatInterval, "5s")
	localThreshold := getTimeout(c.LocalThreshold, "5s")
	maxConnIdleTime := getTimeout(c.MaxConnIdleTime, "5s")
	serverSelectionTimeout := getTimeout(c.ServerSelectionTimeout, "5s")
	socketTimeout := getTimeout(c.SocketTimeout, "60s")

	config := &mongodb.Config{
		AppName:                c.AppName,
		Host:                   strings.ReplaceAll(c.Host, " ", ""),
		DatabaseName:           c.DbName,
		UserName:               c.User,
		Password:               c.Pwd,
		WriteConcern:           c.WriteConcern,
		ReadConcern:            c.ReadConcern,
		Direct:                 c.Direct,
		ReplicaSet:             c.ReplicaSet,
		MaxPoolSize:            getInt(c.MaxPoolSize, 20),
		OperationTimeout:       operationTimeout,
		ConnectTimeout:         connectTimeout,
		HeartbeatInterval:      heartbeatInterval,
		LocalThreshold:         localThreshold,
		MaxConnIdleTime:        maxConnIdleTime,
		ServerSelectionTimeout: serverSelectionTimeout,
		SocketTimeout:          socketTimeout,
		AuthMechanism:          c.AuthMechanism,
		AuthSource:             c.AutoSource,
	}
	return config
}

func getInt(v *uint64, def uint64) uint64 {
	if v == nil {
		return def
	}
	return *v
}

func getTimeout(val string, def string) time.Duration {
	val = strings.Trim(val, " ")
	if val == "" {
		val = def
	}
	v, _ := time.ParseDuration(val)
	return v
}

func getPInt(v *uint64) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", intutils.P2Uint64(v))
}
func getPStr(v *string) string {
	if v == nil {
		return "nil"
	}
	return *v
}
