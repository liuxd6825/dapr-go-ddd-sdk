package restapp

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
	"time"
)

type MongoConfig struct {
	DbKey        string
	Host         string  `yaml:"host"`
	Database     string  `yaml:"dbname"`
	UserName     string  `yaml:"user"`
	Password     string  `yaml:"pwd"`
	ReplicaSet   string  `yaml:"replicaSet"`
	WriteConcern string  `yaml:"writeConcern"`
	ReadConcern  string  `yaml:"readConcern"`
	MaxPoolSize  *uint64 `yaml:"maxPoolSize"`

	Direct                 *bool  `json:"direct"`
	LocalThreshold         string `yaml:"localThreshold"`         // 时间长度
	ConnectTimeout         string `yaml:"connectTimeout"`         // 时间长度
	HeartbeatInterval      string `yaml:"heartbeatInterval"`      // 时间长度
	OperationTimeout       string `yaml:"operationTimeout"`       // 时间长度
	MaxConnIdleTime        string `yaml:"maxConnIdleTime"`        // 时间长度
	ServerSelectionTimeout string `yaml:"serverSelectionTimeout"` // 时间长度
	SocketTimeout          string `yaml:"socketTimeout"`          // 时间长度
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

func initMongo(appName string, appMongoConfigs map[string]*MongoConfig) error {
	if appMongoConfigs == nil {
		return nil
	}

	if _initMongo {
		return nil
	}
	_initMongo = true

	if err := assert.NotNil(appMongoConfigs, assert.NewOptions("appMongoConfig is nil")); err != nil {
		panic(err)
	}

	for k, c := range appMongoConfigs {
		if c.IsEmpty() {
			continue
		}

		operationTimeout := defaultTimeout(c.OperationTimeout, "30s")
		connectTimeout := defaultTimeout(c.ConnectTimeout, "5s")
		heartbeatInterval := defaultTimeout(c.HeartbeatInterval, "5s")
		localThreshold := defaultTimeout(c.LocalThreshold, "5s")
		maxConnIdleTime := defaultTimeout(c.MaxConnIdleTime, "5s")
		serverSelectionTimeout := defaultTimeout(c.ServerSelectionTimeout, "5s")
		socketTimeout := defaultTimeout(c.SocketTimeout, "60s")

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
			OperationTimeout:       operationTimeout,
			ConnectTimeout:         connectTimeout,
			HeartbeatInterval:      heartbeatInterval,
			LocalThreshold:         localThreshold,
			MaxConnIdleTime:        maxConnIdleTime,
			ServerSelectionTimeout: serverSelectionTimeout,
			SocketTimeout:          socketTimeout,
		}
		mongodb, err := ddd_mongodb.NewMongoDB(config, func(opts *options.ClientOptions) error {
			GetLogger().Infof("config mongo  hosts=%v; user=%s; replicasSet=%s; maxPoolSize=%s; connectTimeout=%v; "+
				"socketTimeout=%v; serverSelectionTimeout=%v; maxConnIdleTime=%v; operationTimeout=%v; socketTimeout=%v ",
				opts.Hosts, opts.Auth.Username, pstr(opts.ReplicaSet), pint(opts.MaxPoolSize), connectTimeout,
				socketTimeout, serverSelectionTimeout, maxConnIdleTime, operationTimeout, socketTimeout)
			opts.Monitor = newMongoMonitor()
			opts.ServerMonitor = newMongoServerMonitor()
			return nil
		})
		if err != nil {
			logs.Errorf(context.Background(), "", nil, "连接mongo失败, error:%s", config.Host, err.Error())
			os.Exit(1)
		}
		dbKey := strings.ToLower(k)
		c.DbKey = dbKey
		_mongoDbs[dbKey] = mongodb
		_mongoDefault = mongodb
	}
	if len(_mongoDbs) > 1 {
		_mongoDefault = nil
	}
	return nil
}

func newMongoMonitor() *event.CommandMonitor {
	monitor := &event.CommandMonitor{}
	monitor.Failed = func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
		logs.Errorf(context.Background(), "", nil, "Mongodb CommandMonitor Failed %s", func() any {
			if text, err := jsonutils.Marshal(failedEvent); err != nil {
				return err.Error()
			} else {
				return text
			}
		})
	}
	/*
		monitor.Succeeded = func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
			logs.Infof(context.Background(), "", nil, "Mongodb CommandMonitor Succeeded %s", func() any {
				if text, err := jsonutils.Marshal(succeededEvent); err != nil {
					return err.Error()
				} else {
					return text
				}
			})
		}
		monitor.Started = func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			logs.Infof(context.Background(), "", nil, "Mongodb CommandMonitor Started %s", func() any {
				if text, err := jsonutils.Marshal(startedEvent); err != nil {
					return err.Error()
				} else {
					return text
				}
			})
		}
	*/
	return monitor
}

func newMongoServerMonitor() *event.ServerMonitor {
	monitor := &event.ServerMonitor{}
	monitor.ServerClosed = func(event *event.ServerClosedEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor ServerClosed Address=%s, TopologyID=%v  ", event.Address, event.TopologyID)
	}
	monitor.ServerHeartbeatFailed = func(event *event.ServerHeartbeatFailedEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor ServerHeartbeatFailed Failure=%s ", event.Failure.Error())
	}
	monitor.ServerDescriptionChanged = func(event *event.ServerDescriptionChangedEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor ServerDescriptionChanged Address=%s, NewDescription=%v  ", event.Address, event.NewDescription)
	}
	monitor.ServerOpening = func(event *event.ServerOpeningEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor ServerOpening Address=%s, TopologyID=%v  ", event.Address, event.TopologyID)
	}
	monitor.TopologyClosed = func(event *event.TopologyClosedEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor TopologyClosed TopologyID=%v ", event.TopologyID)
	}
	monitor.TopologyDescriptionChanged = func(event *event.TopologyDescriptionChangedEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor TopologyDescriptionChanged TopologyID=%s, NewDescription=%v  ", event.TopologyID, event.NewDescription)
	}
	monitor.TopologyOpening = func(event *event.TopologyOpeningEvent) {
		logs.Infof(context.Background(), "", nil, "Mongodb ServerMonitor TopologyOpening TopologyID=%v  ", event.TopologyID)
	}
	return monitor
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

func defaultTimeout(val string, def string) time.Duration {
	val = strings.Trim(val, " ")
	if val == "" {
		val = def
	}
	v, _ := time.ParseDuration(val)
	return v
}
