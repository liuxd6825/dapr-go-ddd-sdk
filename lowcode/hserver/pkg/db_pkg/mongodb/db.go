package mongodb

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type MongoDB struct {
	cfg common.IEnvConfig
}

type DB struct {
	mongodb *ddd_mongodb.MongoDB
	cfg     common.IEnvConfig
}

func New(cfg common.IEnvConfig) *MongoDB {
	return &MongoDB{cfg: cfg}
}

func (d *MongoDB) NewDao(dbName, tableName string) *Dao {
	mongoBb, ok := restapp.GetMongoByKey(dbName)
	if !ok {
		panic(fmt.Sprintf("%s db not found", dbName))
	}
	db := &DB{cfg: d.cfg, mongodb: mongoBb}
	return NewDao(db, tableName, &ModelOptions{MongoDB: mongoBb})
}

/*func (d *DB) Get(name string) error {
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
*/
