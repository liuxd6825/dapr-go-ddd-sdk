package restapp

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

type MongoConfig struct {
	DbKey                  string
	Host                   string  `yaml:"host"`
	Database               string  `yaml:"dbname"`
	UserName               string  `yaml:"user"`
	Password               string  `yaml:"pwd"`
	ReplicaSet             string  `yaml:"replicaSet"`
	WriteConcern           string  `yaml:"writeConcern"`
	ReadConcern            string  `yaml:"readConcern"`
	OperationTimeout       *uint64 `yaml:"operationTimeout"` //单位秒
	MaxPoolSize            *uint64 `yaml:"maxPoolSize"`
	ConnectTimeout         *uint64 `yaml:"connectTimeout"` //单位秒
	HeartbeatInterval      *uint64 `yaml:"heartbeatInterval"`
	LocalThreshold         *uint64 `yaml:"localThreshold"`
	MaxConnIdleTime        *uint64 `yaml:"maxConnIdleTime"`
	ServerSelectionTimeout *uint64 `yaml:"serverSelectionTimeout"`
	SocketTimeout          *uint64 `yaml:"socketTimeout"`
	Direct                 *bool   `json:"direct"`
}

var _mongoDbs map[string]*ddd_mongodb.MongoDB
var _mongoDefault *ddd_mongodb.MongoDB
var _initMongo = false

func (m MongoConfig) IsEmpty() bool {
	if m.Host == "" && m.Database == "" && m.Password == "" && m.UserName == "" {
		return true
	}
	return false
}

func init() {
	_mongoDbs = make(map[string]*ddd_mongodb.MongoDB)
}

func initMongo(appName string, appMongoConfigs map[string]*MongoConfig) {
	if _initMongo {
		return
	}
	_initMongo = true
	if err := assert.NotNil(appMongoConfigs, assert.NewOptions("appMongoConfig is nil")); err != nil {
		panic(err)
	}

	for k, c := range appMongoConfigs {
		if c.IsEmpty() {
			continue
		}
		config := &ddd_mongodb.Config{
			AppName:                appName,
			Host:                   strings.ReplaceAll(c.Host, " ", ""),
			DatabaseName:           c.Database,
			UserName:               c.UserName,
			Password:               c.Password,
			WriteConcern:           c.WriteConcern,
			ReadConcern:            c.ReadConcern,
			Direct:                 c.Direct,
			ReplicaSet:             c.ReplicaSet,
			MaxPoolSize:            defaultInt(c.MaxPoolSize, 20),
			OperationTimeout:       defaultTimeout(c.OperationTimeout, 30),
			ConnectTimeout:         defaultTimeout(c.ConnectTimeout, 5),
			HeartbeatInterval:      defaultTimeout(c.HeartbeatInterval, 5),
			LocalThreshold:         defaultTimeout(c.LocalThreshold, 5),
			MaxConnIdleTime:        defaultTimeout(c.MaxConnIdleTime, 5),
			ServerSelectionTimeout: defaultTimeout(c.ServerSelectionTimeout, 5),
			SocketTimeout:          defaultTimeout(c.SocketTimeout, 60),
		}
		mongodb, err := ddd_mongodb.NewMongoDB(config, func(opts *options.ClientOptions) error {
			_logger.Infof("config mongo  hosts=%v; user=%s; replicasSet=%s; maxPoolSize=%s; connectTimeout=%s; "+
				"socketTimeout=%s; serverSelectionTimeout=%s; maxConnIdleTime=%s; operationTimeout=%v;timeUnit=second; ",
				opts.Hosts, opts.Auth.Username, pstr(opts.ReplicaSet), pint(opts.MaxPoolSize), pseconds(opts.ConnectTimeout),
				pseconds(opts.SocketTimeout), pseconds(opts.ServerSelectionTimeout), pseconds(opts.MaxConnIdleTime), config.OperationTimeout)
			return nil
		})
		if err != nil {
			panic(err)
		}
		dbKey := strings.ToLower(k)
		c.DbKey = dbKey
		_mongoDbs[dbKey] = mongodb
		_mongoDefault = mongodb
	}
	if len(_mongoDbs) > 1 {
		_mongoDefault = nil
	}
}

func GetMongoDB() *ddd_mongodb.MongoDB {
	return _mongoDefault
}

func GetMongoByKey(dbKey string) (*ddd_mongodb.MongoDB, bool) {
	d, ok := _mongoDbs[strings.ToLower(dbKey)]
	return d, ok
}

func CloseMongoDB(ctx context.Context) error {
	c := func(d *ddd_mongodb.MongoDB) (err error) {
		defer func() {
			err = errors.GetRecoverError(err, recover())
		}()
		return d.Close(ctx)
	}
	for _, d := range _mongoDbs {
		_ = c(d)
	}
	return nil
}

func pseconds(v *time.Duration) string {
	if v == nil {
		return "null"
	}
	return fmt.Sprintf("%v", v.Seconds())
}

func pint(v *uint64) string {
	if v == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", intutils.P2Uint64(v))
}
func pstr(v *string) string {
	if v == nil {
		return "nil"
	}
	return *v
}

func defaultInt(v *uint64, def uint64) uint64 {
	if v == nil {
		return def
	}
	return *v
}

func defaultTimeout(v *uint64, def uint64) time.Duration {
	if v == nil {
		return time.Duration(def) * time.Second
	}
	return time.Duration(*v) * time.Second
}
