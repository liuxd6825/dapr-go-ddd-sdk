package restapp

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"strings"
)

type MongoConfig struct {
	Host         string `yaml:"host"`
	Database     string `yaml:"dbname"`
	UserName     string `yaml:"user"`
	Password     string `yaml:"pwd"`
	MaxPoolSize  uint64 `yaml:"maxPoolSize"`
	ReplicaSet   string `yaml:"replicaSet"`
	WriteConcern string `yaml:"writeConcern"`
	ReadConcern  string `yaml:"readConcern"`
}

func (m MongoConfig) IsEmpty() bool {
	if m.Host == "" && m.Database == "" && m.Password == "" && m.UserName == "" {
		return true
	}
	return false
}

var _mongoDbs map[string]*ddd_mongodb.MongoDB
var _mongoDefault *ddd_mongodb.MongoDB

func initMongo(appMongoConfigs map[string]*MongoConfig) {
	if err := assert.NotNil(appMongoConfigs, assert.NewOptions("appMongoConfig is nil")); err != nil {
		panic(err)
	}

	for k, c := range appMongoConfigs {
		if c.IsEmpty() {
			continue
		}
		config := &ddd_mongodb.Config{
			Host:         c.Host,
			DatabaseName: c.Database,
			UserName:     c.UserName,
			Password:     c.Password,
			MaxPoolSize:  c.MaxPoolSize,
			ReplicaSet:   c.ReplicaSet,
			WriteConcern: c.WriteConcern,
			ReadConcern:  c.ReadConcern,
		}
		mongodb, err := ddd_mongodb.NewMongoDB(config)
		if err != nil {
			panic(err)
		}
		_mongoDbs[strings.ToLower(k)] = mongodb
		_mongoDefault = mongodb
	}
	if len(_mongoDbs) > 1 {
		_mongoDefault = nil
	}
}

func GetMongoDB() *ddd_mongodb.MongoDB {
	return _mongoDefault
}

/*
func NewSession(isWrite bool) ddd_repository.Session {
	return ddd_mongodb.NewSession(isWrite, _mongoDefault)
}
*/
func GetMongoByKey(dbKey string) (*ddd_mongodb.MongoDB, bool) {
	d, ok := _mongoDbs[strings.ToLower(dbKey)]
	return d, ok
}

func CloseMongoDB(ctx context.Context) error {
	c := func(d *ddd_mongodb.MongoDB) (err error) {
		defer func() {
			err = errors.GetRecoverError(recover())
		}()
		return d.Close(ctx)
	}
	for _, d := range _mongoDbs {
		_ = c(d)
	}
	return nil
}
