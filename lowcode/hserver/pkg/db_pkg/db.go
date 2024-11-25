package db_pkg

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type DB struct {
	mongodb *ddd_mongodb.MongoDB
	models  map[string]*Model
	cfg     common.IEnvConfig
}

func New(cfg common.IEnvConfig) *DB {
	return &DB{models: map[string]*Model{}, cfg: cfg}
}

func newDb(mongodb *ddd_mongodb.MongoDB) *DB {
	return &DB{mongodb: mongodb, models: map[string]*Model{}, cfg: nil}
}

func (d *DB) Get(name string) error {
	db, ok := restapp.GetMongoByKey(name)
	if !ok {
		panic(fmt.Sprintf("%s db not found", name))
	}
	if db != nil {
		d.mongodb = db
	}
	return nil
}

func (d *DB) Open(cfg restapp.MongoConfig) error {
	config := restapp.NewDddMongodbConfig(&cfg)
	mongodb, err := ddd_mongodb.NewMongoDB(config, nil)
	if mongodb != nil {
		d.mongodb = mongodb
	}
	return err
}

func (d *DB) Model(name string) *Model {
	return NewModel(d, name, &ModelOptions{MongoDB: d.mongodb})
}
