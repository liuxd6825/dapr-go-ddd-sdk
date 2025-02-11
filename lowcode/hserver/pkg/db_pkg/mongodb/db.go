package mongodb

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
)

type MongoDB struct {
	server element.Server
}

type DB struct {
	mongodb *ddd_mongodb.MongoDB
	cfg     common.IEnvConfig
}

type NewDaoOptions struct {
	DbName     string
	TableName  string
	IsPubEvent bool
	AggField   string
}

func New(server element.Server) *MongoDB {
	return &MongoDB{server: server}
}

func (d *MongoDB) NewDao(opts *NewDaoOptions) *Dao {
	mongoBb, ok := restapp.GetMongoByKey(opts.DbName)
	if !ok {
		panic(fmt.Sprintf("%s db not found", opts.DbName))
	}
	db := &DB{cfg: d.server.GetEnvCfg(), mongodb: mongoBb}

	daoOpts := &DaoOptions{
		DB:         db,
		TableName:  opts.TableName,
		MongoDB:    db.mongodb,
		AggField:   opts.AggField,
		IsPubEvent: opts.IsPubEvent,
		Server:     d.server,
	}
	return NewDao(daoOpts)
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
