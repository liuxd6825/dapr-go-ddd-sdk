package db

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type DB struct {
	mongodb *ddd_mongodb.MongoDB
	models  map[string]*Model
}

var db = NewDB()

func NewDB() *DB {
	return &DB{models: map[string]*Model{}}
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
